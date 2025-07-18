// internal/handler/author.go
package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/iliyamo/go-learning/internal/model"
	"github.com/iliyamo/go-learning/internal/repository"
	"github.com/labstack/echo/v4"
)

// AuthorRequest ساختار ورودی برای ایجاد یا ویرایش نویسنده است
// این داده‌ها از سمت کاربر دریافت می‌شود
// و شامل نام، بیوگرافی و تاریخ تولد نویسنده است.
type AuthorRequest struct {
	Name      string `json:"name"`
	Biography string `json:"biography"`
	BirthDate string `json:"birth_date"`
}

// CreateAuthor ایجاد نویسنده جدید
func CreateAuthor(c echo.Context) error {
	repo := c.Get("author_repo").(*repository.AuthorRepository)
	req := new(AuthorRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "درخواست نامعتبر است"})
	}

	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "فرمت تاریخ تولد نادرست است، از YYYY-MM-DD استفاده کنید"})
	}

	exists, err := repo.Exists(req.Name, birthDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "خطا در بررسی نویسنده"})
	}
	if exists {
		return c.JSON(http.StatusConflict, echo.Map{"error": "نویسنده قبلاً ثبت شده است"})
	}

	author := &model.Author{
		Name:      req.Name,
		Biography: req.Biography,
		BirthDate: birthDate,
		CreatedAt: time.Now(),
	}

	if err := repo.CreateAuthor(author); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "ثبت نویسنده با خطا مواجه شد"})
	}

	return c.JSON(http.StatusCreated, author)
}

// GetAllAuthors دریافت لیست همه نویسنده‌ها
func GetAllAuthors(c echo.Context) error {
	repo := c.Get("author_repo").(*repository.AuthorRepository)
	authors, err := repo.GetAllAuthors()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "دریافت نویسنده‌ها با خطا مواجه شد"})
	}
	return c.JSON(http.StatusOK, authors)
}

// GetAuthorByID دریافت اطلاعات نویسنده با شناسه مشخص
func GetAuthorByID(c echo.Context) error {
	repo := c.Get("author_repo").(*repository.AuthorRepository)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "شناسه نامعتبر است"})
	}
	author, err := repo.GetAuthorByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "خطا در دریافت نویسنده"})
	}
	if author == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "نویسنده یافت نشد"})
	}
	return c.JSON(http.StatusOK, author)
}

// UpdateAuthor بروزرسانی اطلاعات نویسنده
func UpdateAuthor(c echo.Context) error {
	repo := c.Get("author_repo").(*repository.AuthorRepository)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "شناسه نامعتبر است"})
	}

	req := new(AuthorRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "درخواست نامعتبر است"})
	}

	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "فرمت تاریخ تولد نادرست است"})
	}

	author := &model.Author{
		ID:        id,
		Name:      req.Name,
		Biography: req.Biography,
		BirthDate: birthDate,
	}

	updated, err := repo.UpdateAuthor(author)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "خطا در بروزرسانی نویسنده"})
	}
	if !updated {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "نویسنده یافت نشد"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "نویسنده با موفقیت بروزرسانی شد"})
}

// DeleteAuthor حذف نویسنده بر اساس شناسه
func DeleteAuthor(c echo.Context) error {
	repo := c.Get("author_repo").(*repository.AuthorRepository)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "شناسه نامعتبر است"})
	}

	deleted, err := repo.DeleteAuthor(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "خطا در حذف نویسنده"})
	}
	if !deleted {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "نویسنده یافت نشد"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "نویسنده با موفقیت حذف شد"})
}
