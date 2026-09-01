package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/spikycham/feedme/internal/constant"
)

type FoodRepository struct {
	db *sql.DB
}

func NewFoodRepository(db *sql.DB) *FoodRepository {
	return &FoodRepository{db}
}

type (
	FoodCategory int

	Food struct {
		ID           int
		FoodID       string
		Name         string
		Detail       string
		Prize        float32
		Rate         float32
		RequiredTime int64
		SoldCount    int
		ImageURIs    []string
		Category     FoodCategory // 0 staple food, 1 vegetable, 2 meat, 3 seafood, 4 soup, 5 dessert, 6 drink, 7 other
		CreatedAt    int64
		DeletedAt    int64
	}
)

// PERF: i guess the dish types could be mixed, but i would just
// skip this feature and provide only one type to each dish.
const (
	FoodCategoryStaple FoodCategory = iota
	FoodCategoryVegetable
	FoodCategoryMeat
	FoodCategorySeafood
	FoodCategorySoup
	FoodCategoryDessert
	FoodCategoryDrink
	FoodCategoryOther
)

func (r *FoodRepository) SelectAllFoods(ctx context.Context) ([]Food, error) {
	foods := make([]Food, 0)

	rows, err := r.db.QueryContext(ctx, "SELECT food_id, name, detail, prize, rate, required_time, image_uris, category, created_at, deleted_at FROM foods")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var food Food
		var imgUris string

		if err := rows.Scan(
			&food.FoodID,
			&food.Name,
			&food.Detail,
			&food.Prize,
			&food.Rate,
			&food.RequiredTime,
			&imgUris,
			&food.Category,
			&food.CreatedAt,
			&food.DeletedAt,
		); err != nil {
			return nil, err
		}

		food.ID = -1
		food.ImageURIs = strings.Split(imgUris[1:len(imgUris)-1], ",")

		foods = append(foods, food)
	}

	return foods, nil
}

type InsertFoodParams struct {
	FoodID       string
	Name         string
	Detail       string
	Prize        float32
	Rate         float32
	RequiredTime int64
	ImageURIs    []string
	Category     FoodCategory // 0 staple food, 1 vegetable, 2 meat, 3 seafood, 4 soup, 5 dessert, 6 drink, 7 other
}

func (r *FoodRepository) InsertFood(ctx context.Context, p *InsertFoodParams) error {
	if _, err := r.db.ExecContext(
		ctx, `
		INSERT INTO foods (
			food_id, name, detail, prize, rate, required_time, image_uris, category
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		p.FoodID,
		p.Name,
		p.Detail,
		p.Prize,
		p.Rate,
		p.RequiredTime,
		fmt.Sprintf("[%s]", strings.Join(p.ImageURIs, ", ")),
		p.Category,
	); err != nil {
		return err
	}
	return nil
}

type UpdateFoodParams struct {
	Name         *string
	Detail       *string
	Prize        *float32
	Rate         *float32
	RequiredTime *int64
	ImageURIs    *[]string
	Category     *FoodCategory // 0 staple food, 1 vegetable, 2 meat, 3 seafood, 4 soup, 5 dessert, 6 drink, 7 other
}

func (r *FoodRepository) UpdateFoodByFoodID(ctx context.Context, foodId string, p *UpdateFoodParams) error {
	var sets []string
	var args []any

	if p.Name != nil {
		sets = append(sets, "name = ?")
		args = append(args, *p.Name)
	}
	if p.Detail != nil {
		sets = append(sets, "detail = ?")
		args = append(args, *p.Detail)
	}
	if p.Prize != nil {
		sets = append(sets, "prize = ?")
		args = append(args, *p.Prize)
	}
	if p.Rate != nil {
		sets = append(sets, "rate = ?")
		args = append(args, *p.Rate)
	}
	if p.RequiredTime != nil {
		sets = append(sets, "required_time = ?")
		args = append(args, *p.RequiredTime)
	}
	if p.ImageURIs != nil {
		sets = append(sets, "image_uris = ?")
		args = append(args, fmt.Sprintf("[%s]", strings.Join(*p.ImageURIs, ", ")))
	}
	if p.Category != nil {
		sets = append(sets, "category = ?")
		args = append(args, *p.Category)
	}

	if len(sets) == 0 {
		return nil
	}

	args = append(args, foodId)
	query := fmt.Sprintf("UPDATE foods SET %s WHERE food_id = ?", strings.Join(sets, ", "))

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return constant.NoAffectedRows
	}

	return nil
}

func (r *FoodRepository) DeleteFoodByFoodID(ctx context.Context, foodId string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM foods WHERE food_id = ?", foodId)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return constant.NoAffectedRows
	}

	return nil
}
