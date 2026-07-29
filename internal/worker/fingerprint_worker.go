package worker

import (
	"log"
	"time"

	"universev2-backend/internal/model"
	"universev2-backend/internal/repository"
	"universev2-backend/pkg/solutionx100c"
)

type FingerprintWorker struct {
	fpRepo   *repository.FingerprintRepo
	attRepo  *repository.AttendanceRepo
	ticker   *time.Ticker
	stopChan chan struct{}
}

func NewFingerprintWorker(fpRepo *repository.FingerprintRepo, attRepo *repository.AttendanceRepo) *FingerprintWorker {
	return &FingerprintWorker{
		fpRepo:   fpRepo,
		attRepo:  attRepo,
		stopChan: make(chan struct{}),
	}
}

func (w *FingerprintWorker) Start(interval time.Duration) {
	if interval <= 0 {
		interval = 60 * time.Second
	}
	w.ticker = time.NewTicker(interval)

	go func() {
		log.Printf("[FingerprintWorker] Worker started. Polling every %v...", interval)
		// Run immediate initial sync
		w.SyncAllDevices()

		for {
			select {
			case <-w.ticker.C:
				w.SyncAllDevices()
			case <-w.stopChan:
				w.ticker.Stop()
				log.Println("[FingerprintWorker] Worker stopped.")
				return
			}
		}
	}()
}

func (w *FingerprintWorker) Stop() {
	close(w.stopChan)
}

func (w *FingerprintWorker) SyncAllDevices() int {
	devices, err := w.fpRepo.GetActiveDevices()
	if err != nil {
		log.Printf("[FingerprintWorker] Error fetching active devices from DB: %v", err)
		return 0
	}

	totalSynced := 0
	for _, dev := range devices {
		syncedCount, online := w.SyncDevice(&dev)
		_ = w.fpRepo.UpdateSyncStatus(dev.ID, online)
		totalSynced += syncedCount
	}

	if totalSynced > 0 {
		log.Printf("[FingerprintWorker] Successfully synced %d attendance logs across devices.", totalSynced)
	}

	return totalSynced
}

func (w *FingerprintWorker) SyncDevice(dev *model.FingerprintDevice) (int, bool) {
	client := solutionx100c.NewClient(dev.IPAddress, dev.Port, dev.ComKey)

	records, err := client.FetchAttLog()
	if err != nil {
		// Log connection failure, device marked offline
		return 0, false
	}

	synced := 0
	for _, rec := range records {
		machineLabel := dev.Name + " (" + dev.Code + ")"
		if dev.Location != "" {
			machineLabel = dev.Name + " · " + dev.Location
		}

		_, err := w.attRepo.RecordScanWithTimestamp(rec.NIK, machineLabel, rec.Timestamp)
		if err == nil {
			synced++
		}
	}

	// Flush machine memory after reading
	if len(records) > 0 {
		_ = client.ClearAttLog()
	}

	return synced, true
}
