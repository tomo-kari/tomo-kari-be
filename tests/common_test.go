package tests

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"os"
	"testing"
	v1 "tomokari/internal/controller/http/v1"
	"tomokari/internal/entity"
	"tomokari/internal/usecase"
	"tomokari/internal/usecase/repo"
	"tomokari/pkg/logger"
	pgDb "tomokari/pkg/postgres"
)

var pg *pgDb.Postgres
var userRoute *v1.UserRoutes

func TestMain(m *testing.M) {
	//var err error
	//err = godotenv.Load(os.ExpandEnv("../../.env"))
	//if err != nil {
	//    log.Fatalf("Error getting env %v\n", err)
	//}
	Database()
	userUseCase := usecase.NewUserUseCase(
		repo.NewUserRepo(pg),
		repo.NewTOSRepo(pg),
	)
	l := logger.New("debug")
	v := validator.New()
	userRoute = &v1.UserRoutes{
		U: userUseCase,
		L: l,
		V: v,
	}

	os.Exit(m.Run())
}

func Database() {
	dbUrl := "postgres://postgres:postgrespw@localhost:5432/tomokari-test?sslmode=disable"
	//db, err := sql.Open("postgres", dbUrl)
	//if err != nil {
	//    fmt.Printf("Cannot connect to %s database\n", "postgres")
	//    log.Fatal("This is the error:", err)
	//} else {
	//    fmt.Printf("We are connected to the %s database\n", "postgres")
	//}
	//driver, err := postgres.WithInstance(db, &postgres.Config{})
	//if err != nil {
	//    log.Fatal("Error getting driver", err)
	//}
	//m, err := migrate.NewWithDatabaseInstance(
	//    "migrations",
	//    "tomokari-test",
	//    driver,
	//)
	//if err != nil {
	//    log.Fatal("Error getting migrator: ", err)
	//}
	//err = m.Up()
	//if err != nil {
	//    fmt.Printf("Cannot connect to %s database\n, error %v", "postgres", err)
	//    return
	//}
	var err error
	pg, err = pgDb.New(dbUrl, pgDb.MaxPoolSize(3))
	if err != nil {
		fmt.Printf("Cannot connect to %s database\n", "postgres")
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database\n", "postgres")
	}
}

func deleteData(tableName string) error {
	query, args, err := pg.Builder.Delete(tableName).ToSql()
	if err != nil {
		return err
	}
	res, err := pg.Pool.Exec(context.Background(), query, args...)
	if err != nil {
		return err
	}
	log.Printf("%d rows deleted from %s table", res.RowsAffected(), tableName)
	return nil
}

func refreshUserTable() error {
	query, args, err := pg.Builder.Delete("users").ToSql()
	if err != nil {
		return err
	}
	res, err := pg.Pool.Exec(context.Background(), query, args...)
	if err != nil {
		return err
	}
	log.Printf("%d rows deleted from users table", res.RowsAffected())
	return nil
}

var tos = entity.TermsOfService{
	ID:      1,
	Content: "test",
}

func seedOneTos() {
	err := deleteData((&entity.TermsOfService{}).Table())
	if err != nil {
		fmt.Printf("Error refreshing user table %v\n", err)
	}
	query, args, err := pg.Builder.
		Insert((&entity.TermsOfService{}).Table()).
		Columns("id, content").
		//Values(user.Email, user.Phone, user.Password, user.Role, user.DateOfBirth, user.Description).
		Values(tos.ID, tos.Content).
		ToSql()
	if err != nil {
		fmt.Println(err)
	}
	ctx := context.Background()
	pg.Pool.QueryRow(ctx, query, args...)
}
