package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/spikycham/feedme/internal/constant"
	"github.com/spikycham/feedme/internal/model"
)

type FoodRepository struct {
	db *sql.DB
}

func NewFoodRepository(db *sql.DB) *FoodRepository {
	return &FoodRepository{db}
}

func (r *FoodRepository) SelectAllFoods(ctx context.Context) ([]model.FoodDetail, error) {
	foodDetails := make([]*model.FoodDetail, 0)
	foodMap := make(map[string]*model.FoodDetail)

	foods, err := r.db.QueryContext(
		ctx,
		"SELECT food_id, name, detail, prize, rate, required_time, image_uris, ingredients, category, created_at, deleted_at FROM foods",
	)
	if err != nil {
		return nil, err
	}
	defer foods.Close()

	for foods.Next() {
		detail := &model.FoodDetail{}
		var imgUris string
		var ingredients string

		if err := foods.Scan(
			&detail.Food.FoodID,
			&detail.Food.Name,
			&detail.Food.Detail,
			&detail.Food.Prize,
			&detail.Food.Rate,
			&detail.Food.RequiredTime,
			&imgUris,
			&ingredients,
			&detail.Food.Category,
			&detail.Food.CreatedAt,
			&detail.Food.DeletedAt,
		); err != nil {
			return nil, err
		}

		imgUrisContent := strings.Split(imgUris[1:len(imgUris)-1], ",")
		if imgUrisContent[0] == "" {
			detail.Food.ImageURIs = []string{}
		} else {
			detail.Food.ImageURIs = imgUrisContent
		}

		ingredientsContent := strings.Split(ingredients[1:len(ingredients)-1], ",")
		if ingredientsContent[0] == "" {
			detail.Food.Ingredients = []string{}
		} else {
			detail.Food.Ingredients = ingredientsContent
		}

		foodDetails = append(foodDetails, detail)
		foodMap[detail.Food.FoodID] = detail
	}

	steps, err := r.db.QueryContext(
		ctx,
		"SELECT food_id, sort, detail FROM food_steps ORDER BY food_id, sort",
	)
	if err != nil {
		return nil, err
	}
	defer steps.Close()

	for steps.Next() {
		var foodId string
		var step model.FoodStep
		if err := steps.Scan(&foodId, &step.Sort, &step.Detail); err != nil {
			return nil, err
		}

		if detail, ok := foodMap[foodId]; ok {
			detail.Steps = append(detail.Steps, step)
		}
	}

	comments, err := r.db.QueryContext(
		ctx,
		"SELECT food_id, comment_id, detail, created_at, deleted_at FROM food_comments ORDER BY food_id",
	)
	if err != nil {
		return nil, err
	}

	for comments.Next() {
		var foodId string
		var comment model.FoodComment
		if err := comments.Scan(&foodId, &comment.CommentID, &comment.Detail, &comment.CreatedAt, &comment.DeletedAt); err != nil {
			return nil, err
		}

		if detail, ok := foodMap[foodId]; ok {
			detail.Comments = append(detail.Comments, comment)
		}
	}

	result := make([]model.FoodDetail, len(foodDetails))
	for i, detail := range foodDetails {
		result[i] = *detail
	}

	return result, nil
}

type InsertFoodParams struct {
	FoodID       string
	Name         string
	Detail       string
	Prize        float32
	Rate         float32
	RequiredTime int64
	ImageURIs    []string
	Ingredients  []string
	Category     model.FoodCategory // 0 staple food, 1 vegetable, 2 meat, 3 seafood, 4 soup, 5 dessert, 6 drink, 7 other
	Steps        []model.FoodStep
}

func (r *FoodRepository) InsertFood(ctx context.Context, p *InsertFoodParams) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(
		ctx, `
		INSERT INTO foods (
			food_id, name, detail, prize, rate, required_time, image_uris, ingredients, category
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		p.FoodID,
		p.Name,
		p.Detail,
		p.Prize,
		p.Rate,
		p.RequiredTime,
		fmt.Sprintf("[%s]", strings.Join(p.ImageURIs, ",")),
		fmt.Sprintf("[%s]", strings.Join(p.Ingredients, ",")),
		p.Category,
	); err != nil {
		return err
	}

	for _, s := range p.Steps {
		if _, err := tx.ExecContext(
			ctx,
			`
			INSERT INTO food_steps (
				food_id, sort, detail
			) VALUES (?, ?, ?)
			`,
			s.FoodID,
			s.Sort,
			s.Detail,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

type UpdateFoodParams struct {
	Name         *string
	Detail       *string
	Prize        *float32
	Rate         *float32
	RequiredTime *int64
	ImageURIs    []string
	Ingredients  []string            `json:"ingredients"`
	Category     *model.FoodCategory // 0 staple food, 1 vegetable, 2 meat, 3 seafood, 4 soup, 5 dessert, 6 drink, 7 other
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
		args = append(args, fmt.Sprintf("[%s]", strings.Join(p.ImageURIs, ", ")))
	}
	if p.Ingredients != nil {
		sets = append(sets, "ingredients = ?")
		args = append(args, fmt.Sprintf("[%s]", strings.Join(p.Ingredients, ", ")))
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

func (r *FoodRepository) UpdateFoodStepsByFoodID(ctx context.Context, foodId string, steps []model.FoodStep) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(
		ctx,
		"DELETE FROM food_steps WHERE food_id = ?",
		foodId,
	); err != nil {
		return err
	}

	for _, s := range steps {
		if _, err := tx.ExecContext(
			ctx,
			`
			INSERT INTO food_steps (
				food_id, sort, detail
			) VALUES (?, ?, ?)
			`,
			foodId,
			s.Sort,
			s.Detail,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
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

// Food comment operations.
func (r *FoodRepository) InsertFoodCommentByFoodID(ctx context.Context, commentId, foodId, detail string) error {
	if _, err := r.db.ExecContext(
		ctx,
		"INSERT INTO food_comments (comment_id, food_id, detail) VALUES (?, ?, ?)",
		commentId,
		foodId,
		detail,
	); err != nil {
		return err
	}

	return nil
}

func (r *FoodRepository) DeleteFoodCommentByCommentID(ctx context.Context, commentId string) error {
	res, err := r.db.ExecContext(ctx, "UPDATE food_comments SET deleted_at = (unixepoch()) WHERE comment_id = ? AND deleted_at = -1", commentId)
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
