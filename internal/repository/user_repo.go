package repository

import (
	"gorm.io/gorm"

	"universev2-backend/internal/model"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetAll() ([]model.User, error) {
	var users []model.User
	if err := r.db.Order("id ASC").Find(&users).Error; err != nil {
		return nil, err
	}

	for i, u := range users {
		var urs []model.UserRole
		r.db.Where("user_id = ?", u.ID).Find(&urs)
		var roles []string
		for _, ur := range urs {
			roles = append(roles, ur.RoleID)
		}
		users[i].Roles = roles
	}

	return users, nil
}

func (r *UserRepo) GetByID(id string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	var urs []model.UserRole
	r.db.Where("user_id = ?", id).Find(&urs)
	for _, ur := range urs {
		user.Roles = append(user.Roles, ur.RoleID)
	}
	return &user, nil
}

func (r *UserRepo) GetByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	var urs []model.UserRole
	r.db.Where("user_id = ?", user.ID).Find(&urs)
	for _, ur := range urs {
		user.Roles = append(user.Roles, ur.RoleID)
	}
	return &user, nil
}

func (r *UserRepo) Create(user *model.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	for _, roleID := range user.Roles {
		r.db.Create(&model.UserRole{UserID: user.ID, RoleID: roleID})
	}
	return nil
}

func (r *UserRepo) Update(id string, user *model.User) error {
	if err := r.db.Model(&model.User{}).Where("id = ?", id).Updates(user).Error; err != nil {
		return err
	}
	r.db.Where("user_id = ?", id).Delete(&model.UserRole{})
	for _, roleID := range user.Roles {
		r.db.Create(&model.UserRole{UserID: id, RoleID: roleID})
	}
	return nil
}

func (r *UserRepo) Delete(id string) error {
	r.db.Where("user_id = ?", id).Delete(&model.UserRole{})
	return r.db.Where("id = ?", id).Delete(&model.User{}).Error
}

func (r *UserRepo) ExistsByEmail(email string) bool {
	var count int64
	r.db.Model(&model.User{}).Where("LOWER(email) = LOWER(?)", email).Count(&count)
	return count > 0
}

func (r *UserRepo) ExistsByNIK(nik string) bool {
	var count int64
	r.db.Model(&model.User{}).Where("nik = ?", nik).Count(&count)
	return count > 0
}

func (r *UserRepo) UpdatePassword(id, hash, salt string) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"password_hash": hash,
		"password_salt": salt,
	}).Error
}

func (r *UserRepo) ToggleStatus(id string, active bool) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("is_active", active).Error
}


// RoleRepo

type RoleRepo struct {
	db *gorm.DB
}

func NewRoleRepo(db *gorm.DB) *RoleRepo {
	return &RoleRepo{db: db}
}

func (r *RoleRepo) GetAll() ([]model.Role, error) {
	var roles []model.Role
	if err := r.db.Order("id ASC").Find(&roles).Error; err != nil {
		return nil, err
	}

	for i, role := range roles {
		var perms []model.RolePermission
		r.db.Where("role_id = ?", role.ID).Find(&perms)
		pMap := make(map[string]string)
		for _, p := range perms {
			pMap[p.ModuleName] = p.PermissionLevel
		}
		roles[i].Perms = pMap
	}

	return roles, nil
}

func (r *RoleRepo) Create(role *model.Role) error {
	if err := r.db.Create(role).Error; err != nil {
		return err
	}
	for mod, perm := range role.Perms {
		r.db.Create(&model.RolePermission{RoleID: role.ID, ModuleName: mod, PermissionLevel: perm})
	}
	return nil
}

func (r *RoleRepo) Update(id string, role *model.Role) error {
	if err := r.db.Model(&model.Role{}).Where("id = ?", id).Updates(role).Error; err != nil {
		return err
	}
	r.db.Where("role_id = ?", id).Delete(&model.RolePermission{})
	for mod, perm := range role.Perms {
		r.db.Create(&model.RolePermission{RoleID: id, ModuleName: mod, PermissionLevel: perm})
	}
	return nil
}

func (r *RoleRepo) Delete(id string) error {
	r.db.Where("role_id = ?", id).Delete(&model.RolePermission{})
	return r.db.Where("id = ?", id).Delete(&model.Role{}).Error
}

func (r *RoleRepo) GetByID(id string) (*model.Role, error) {
	var role model.Role
	if err := r.db.Where("id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}
	var perms []model.RolePermission
	r.db.Where("role_id = ?", id).Find(&perms)
	pMap := make(map[string]string)
	for _, p := range perms {
		pMap[p.ModuleName] = p.PermissionLevel
	}
	role.Perms = pMap
	return &role, nil
}

func (r *RoleRepo) CountUsersByRoleID(roleID string) (int64, error) {
	var count int64
	err := r.db.Model(&model.UserRole{}).Where("role_id = ?", roleID).Count(&count).Error
	return count, err
}

func (r *RoleRepo) GetPermissionsForRoles(roleIDs []string) (map[string]string, error) {
	var perms []model.RolePermission
	if err := r.db.Where("role_id IN ?", roleIDs).Find(&perms).Error; err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, p := range perms {
		current, exists := result[p.ModuleName]
		if !exists || permLevel(p.PermissionLevel) > permLevel(current) {
			result[p.ModuleName] = p.PermissionLevel
		}
	}
	return result, nil
}

func permLevel(p string) int {
	switch p {
	case "manage":
		return 2
	case "view":
		return 1
	default:
		return 0
	}
}
