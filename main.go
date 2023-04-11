package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	_ "github.com/pelletier/go-toml/query"
)

type Book struct {
	ID        int    `json:"id"`
	Judul     string `json:"tittle"`
	Pengarang string `json:"author"`
	Deskripsi string `json:"desc"`
}

var (
	db  *sql.DB
	err error
)

func main() {

	db, err = sql.Open("postgres", "host = localhost port = 5432 user=postgres password = 270300 dbname = db-go-sql sslmode=disable")

	if err != nil {
		panic(err)
	}
	db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println(db)

	g := gin.Default()
	g.GET("/book", getAllBook)
	g.GET("/book/:id", getBookById)
	g.POST("/book", addBook)
	g.DELETE("/book/:id", deleteBook)
	g.PUT("/book/:id", updateBook)
	g.Run(":8080")

}

func getAllBook(ctx *gin.Context) {
	query := "select *from book"
	rows, err := db.Query(query)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	books := make([]Book, 0)
	for rows.Next() {
		var book Book
		err = rows.Scan(&book.ID, &book.Judul, &book.Pengarang, &book.Deskripsi)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})

		}
		books = append(books, book)
	}
	ctx.JSON(http.StatusOK, books)

}
func getBookById(ctx *gin.Context) {

	idString := ctx.Param("id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	var byID Book
	query := "select *from book where id_buku = $1 "
	rows := db.QueryRow(query, id)
	err = rows.Scan(&byID.ID, &byID.Judul, &byID.Pengarang, &byID.Deskripsi)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, byID)

}
func addBook(ctx *gin.Context) {
	var newBook Book

	err := ctx.ShouldBindJSON(&newBook)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return

	}
	query := "insert into book (id_buku, judul, pengarang, deskripsi)values ($1, $2, $3, $4) returning *"
	result, err := db.Exec(query, newBook.ID, newBook.Judul, newBook.Pengarang, newBook.Deskripsi)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	} else {
		ctx.Writer.Write([]byte("Created\n"))
	}
	ctx.JSON(http.StatusOK, result)

}
func deleteBook(ctx *gin.Context) {
	idString := ctx.Param("id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	var deleteBook Book

	query := "delete from book where id_buku = $1 returning *"

	row := db.QueryRow(query, id)

	err = row.Scan(&deleteBook.ID, &deleteBook.Judul, &deleteBook.Pengarang,
		&deleteBook.Deskripsi)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	} else {
		ctx.Writer.Write([]byte("Deleted\n"))

	}
	ctx.JSON(http.StatusOK, deleteBook)

}

func updateBook(ctx *gin.Context) {
	idString := ctx.Param("id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	var newUpdateBook Book
	err = ctx.ShouldBindJSON(&newUpdateBook)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	query := "update book set judul = $1, pengarang = $2, deskripsi = $3 where id_buku = $4 returning id_buku"
	row := db.QueryRow(query, newUpdateBook.Judul, newUpdateBook.Pengarang, newUpdateBook.Deskripsi, id)
	err = row.Scan(&newUpdateBook.ID)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	} else {
		ctx.Writer.Write([]byte("Updated\n"))
	}
	ctx.JSON(http.StatusOK, newUpdateBook)

}
