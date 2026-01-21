package repository

import (
	"ramah-disabilitas-be/internal/model"
	"ramah-disabilitas-be/pkg/database"

	"gorm.io/gorm"
)

func CreateCourse(course *model.Course) error {
	return database.DB.Create(course).Error
}

func GetCoursesByTeacherID(teacherID uint64, search string, status string, sort string) ([]model.Course, error) {
	var courses []model.Course
	query := database.DB.Where("teacher_id = ?", teacherID)

	if search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}

	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	if sort == "oldest" {
		query = query.Order("created_at asc")
	} else {
		query = query.Order("created_at desc")
	}

	err := query.Find(&courses).Error
	return courses, err
}

func GetCourseByID(id uint64) (*model.Course, error) {
	var course model.Course
	err := database.DB.Preload("Modules.Materials").First(&course, id).Error
	return &course, err
}

func UpdateCourse(course *model.Course) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Get IDs of ALL modules currently in DB for this course (to detect deletions)
		var oldModuleIDs []uint64
		if err := tx.Model(&model.Module{}).Where("course_id = ?", course.ID).Pluck("id", &oldModuleIDs).Error; err != nil {
			return err
		}

		// Map of old module IDs -> its old material IDs (to detect material deletions)
		oldMaterialsMap := make(map[uint64][]uint64)
		for _, mid := range oldModuleIDs {
			var oldMatIDs []uint64
			if err := tx.Model(&model.Material{}).Where("module_id = ?", mid).Pluck("id", &oldMatIDs).Error; err != nil {
				return err
			}
			oldMaterialsMap[mid] = oldMatIDs
		}

		// 2. Save Everything (Updates & Creates)
		if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(course).Error; err != nil {
			return err
		}

		// 3. Delete Orphans
		var keepModuleIDs []uint64
		for _, m := range course.Modules {
			if m.ID != 0 {
				keepModuleIDs = append(keepModuleIDs, m.ID)
			}
		}

		for _, oldID := range oldModuleIDs {
			found := false
			for _, keepID := range keepModuleIDs {
				if oldID == keepID {
					found = true
					break
				}
			}
			if !found {
				// Delete Module (and sure to delete materials first if no cascade)
				// 1. Get IDs of materials to be deleted
				mIDs := oldMaterialsMap[oldID]
				if len(mIDs) > 0 {
					// Delete dependencies first
					if err := tx.Where("material_id IN ?", mIDs).Delete(&model.MaterialCompletion{}).Error; err != nil {
						return err
					}
					if err := tx.Where("material_id IN ?", mIDs).Delete(&model.SmartFeature{}).Error; err != nil {
						return err
					}
				}

				if err := tx.Where("module_id = ?", oldID).Delete(&model.Material{}).Error; err != nil {
					return err
				}
				if err := tx.Delete(&model.Module{}, oldID).Error; err != nil {
					return err
				}
			} else {
				// Check for orphaned materials in this kept module
				var currentModule *model.Module
				for i := range course.Modules {
					if course.Modules[i].ID == oldID {
						currentModule = &course.Modules[i]
						break
					}
				}

				if currentModule != nil {
					var keepMatIDs []uint64
					for _, mat := range currentModule.Materials {
						if mat.ID != 0 {
							keepMatIDs = append(keepMatIDs, mat.ID)
						}
					}

					oldMatIDs := oldMaterialsMap[oldID]
					for _, oldMatID := range oldMatIDs {
						foundMat := false
						for _, keepMatID := range keepMatIDs {
							if oldMatID == keepMatID {
								foundMat = true
								break
							}
						}
						if !foundMat {
							// Delete dependencies first
							if err := tx.Where("material_id = ?", oldMatID).Delete(&model.MaterialCompletion{}).Error; err != nil {
								return err
							}
							if err := tx.Where("material_id = ?", oldMatID).Delete(&model.SmartFeature{}).Error; err != nil {
								return err
							}

							if err := tx.Delete(&model.Material{}, oldMatID).Error; err != nil {
								return err
							}
						}
					}
				}
			}
		}

		return nil
	})
}

func DeleteCourse(id uint64) error {
	return database.DB.Delete(&model.Course{}, id).Error
}

func GetCourseByClassCode(code string) (*model.Course, error) {
	var course model.Course
	err := database.DB.Where("class_code = ?", code).First(&course).Error
	return &course, err
}

func AddStudentToCourse(courseID, studentID uint64) error {
	course := model.Course{ID: courseID}
	student := model.User{ID: studentID}
	return database.DB.Model(&course).Association("Students").Append(&student)
}

func IsStudentInCourse(courseID, studentID uint64) (bool, error) {
	var count int64
	err := database.DB.Table("course_students").Where("course_id = ? AND user_id = ?", courseID, studentID).Count(&count).Error
	return count > 0, err
}

func GetCoursesByStudentID(studentID uint64) ([]model.Course, error) {
	var courses []model.Course
	err := database.DB.Table("courses").
		Joins("JOIN course_students ON courses.id = course_students.course_id").
		Where("course_students.user_id = ?", studentID).
		Find(&courses).Error

	if err != nil {
		return nil, err
	}

	for i := range courses {
		courses[i].Progress = calculateCourseProgress(courses[i].ID, studentID)
	}

	return courses, nil
}

func calculateCourseProgress(courseID, studentID uint64) float64 {
	var totalMaterials int64
	var completedMaterials int64
	var totalAssignments int64
	var submittedAssignments int64

	// Count Materials
	database.DB.Table("materials").
		Joins("JOIN modules ON materials.module_id = modules.id").
		Where("modules.course_id = ?", courseID).
		Count(&totalMaterials)

	// Count Assignments
	database.DB.Model(&model.Assignment{}).Where("course_id = ?", courseID).Count(&totalAssignments)

	totalItems := totalMaterials + totalAssignments
	if totalItems == 0 {
		return 0
	}

	// Count Completed Materials
	database.DB.Table("material_completions").
		Joins("JOIN materials ON material_completions.material_id = materials.id").
		Joins("JOIN modules ON materials.module_id = modules.id").
		Where("modules.course_id = ? AND material_completions.user_id = ? AND material_completions.completed = ?", courseID, studentID, true).
		Count(&completedMaterials)

	// Count Submitted Assignments
	// Assuming submission implies completion for progress
	database.DB.Table("submissions").
		Joins("JOIN assignments ON submissions.assignment_id = assignments.id").
		Where("assignments.course_id = ? AND submissions.student_id = ?", courseID, studentID).
		Count(&submittedAssignments)

	completedItems := completedMaterials + submittedAssignments

	return (float64(completedItems) / float64(totalItems)) * 100
}

func ToggleMaterialCompletion(userID, materialID uint64) (bool, error) {
	var completion model.MaterialCompletion
	err := database.DB.Where("user_id = ? AND material_id = ?", userID, materialID).First(&completion).Error

	if err == gorm.ErrRecordNotFound {
		// Create as completed
		newCompletion := model.MaterialCompletion{
			UserID:     userID,
			MaterialID: materialID,
			Completed:  true,
		}
		if err := database.DB.Create(&newCompletion).Error; err != nil {
			return false, err
		}
		return true, nil
	} else if err != nil {
		return false, err
	}

	// Toggle
	completion.Completed = !completion.Completed
	if err := database.DB.Save(&completion).Error; err != nil {
		return false, err
	}

	return completion.Completed, nil
}

func GetMaterialCompletionStatus(userID, materialID uint64) bool {
	var count int64
	database.DB.Model(&model.MaterialCompletion{}).
		Where("user_id = ? AND material_id = ? AND completed = ?", userID, materialID, true).
		Count(&count)
	return count > 0
}

func GetCompletedMaterialsMap(courseID, studentID uint64) (map[uint64]bool, error) {
	var completedMaterials []uint64
	// Find all material IDs completed by this user in this course
	err := database.DB.Table("material_completions").
		Joins("JOIN materials ON material_completions.material_id = materials.id").
		Joins("JOIN modules ON materials.module_id = modules.id").
		Where("modules.course_id = ? AND material_completions.user_id = ? AND material_completions.completed = ?", courseID, studentID, true).
		Pluck("material_completions.material_id", &completedMaterials).Error

	if err != nil {
		return nil, err
	}

	result := make(map[uint64]bool)
	for _, id := range completedMaterials {
		result[id] = true
	}
	return result, nil
}

func GetModuleByID(id uint64) (*model.Module, error) {
	var module model.Module
	err := database.DB.First(&module, id).Error
	return &module, err
}

func DeleteModule(id uint64) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// Find materials
		var materialIDs []uint64
		if err := tx.Model(&model.Material{}).Where("module_id = ?", id).Pluck("id", &materialIDs).Error; err != nil {
			return err
		}

		if len(materialIDs) > 0 {
			if err := tx.Where("material_id IN ?", materialIDs).Delete(&model.MaterialCompletion{}).Error; err != nil {
				return err
			}
			if err := tx.Where("material_id IN ?", materialIDs).Delete(&model.SmartFeature{}).Error; err != nil {
				return err
			}
			// Delete materials manually (explicitly)
			if err := tx.Where("module_id = ?", id).Delete(&model.Material{}).Error; err != nil {
				return err
			}
		}

		return tx.Delete(&model.Module{}, id).Error
	})
}

func GetMaterialByID(id uint64) (*model.Material, error) {
	var material model.Material
	err := database.DB.Preload("SmartFeature").First(&material, id).Error
	return &material, err
}

func DeleteMaterial(id uint64) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("material_id = ?", id).Delete(&model.MaterialCompletion{}).Error; err != nil {
			return err
		}
		if err := tx.Where("material_id = ?", id).Delete(&model.SmartFeature{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Material{}, id).Error
	})
}

func CreateMaterial(material *model.Material) error {
	return database.DB.Create(material).Error
}

func UpdateMaterial(material *model.Material) error {
	return database.DB.Save(material).Error
}

func SaveSmartFeature(feature *model.SmartFeature) error {
	return database.DB.Save(feature).Error
}

func GetStudentsByCourseID(courseID uint64) ([]model.User, error) {
	var students []model.User
	err := database.DB.Model(&model.User{}).
		Joins("JOIN course_students ON users.id = course_students.user_id").
		Where("course_students.course_id = ?", courseID).
		Preload("Accessibility").
		Find(&students).Error
	return students, err
}
