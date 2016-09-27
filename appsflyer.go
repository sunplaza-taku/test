package src

import (
//	"ds" //自作のライブラリ群(後述します)
	"html/template"
//	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)
//初期化
func init() {
	http.Handle("/", GetMainEngine())
}

func GetMainEngine() *gin.Engine {
	LoadTemplaytes()
	server := gin.Default()
	server.LoadHTMLGlob("templates/*")
	server.GET("/", Listing)
	server.GET("/get/:id", GetTodo)
	server.POST("/post/", MakeTodo)
	return server
}
//localhost:8080/にアクセス(GETリクエスト)した時の処理。
//Datastoreに存在する全てのTodoを列挙し、そのStringIDを抽出。
//それをHTMLテンプレートに流し込みリンクとして列挙するようにしています。
func Listing(g *gin.Context) {
	keys, err := DataStoreManager.ListingTodoData(g)
	if err != nil {
		g.String(500, "Get Error at Listing %v", err)
		return
	}
	ids := make([]string, len(keys))
	for i, v := range keys {
		ids[i] = v.StringID()
	}
	g.HTML(200, "t.html", gin.H{
		"KeyIDs": ids,
	})
}
//localhost:8080/get/(ID数値)にアクセス(GETリクエスト)した時の処理
//IDを抽出し、それをKeyとしてTodoのEntityを構造体として持ってくる。
//最後にそれをJSONで返している。
func GetTodo(g *gin.Context) {
	ids := g.Param("id")
	todo, err := DataStoreManager.GetTodoList(ids, g)
	if err != nil {
		g.String(500, "Get Error at GetTodo %v", err)
	}
	g.JSON(200, todo)
}

//localhost:8080/post/にPOSTがあった時に行われる処理
//curl -X POST -H 'Content-Type: application/json' -d '{"Title" : "Title1" , "Description" : "いろいろやる"}' localhost:8080/post/
//↑のような例をやると正しくDatastoreに格納される。
func MakeTodo(g *gin.Context) {
	var todo DataStoreManager.Todo
	if g.BindJSON(&todo) == nil {
		if err := DataStoreManager.PutToDoList(todo, g); err != nil {
			g.String(500, "Put Error %v", err)
			return
		}
		g.String(200, "PutSuccess")
	} else {
		g.String(500, "Cant Make JSON")
	}
}
//HTMLテンプレートをロードする。
//データを列挙するページでのみ使用している。
func LoadTemplaytes() {
	baseTemplate, err := template.New("root").Parse("t.html")
	template.Must(todoTemplate, err2)
}
