package network

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
	"todostudy/types"

	"github.com/gorilla/mux"
	"github.com/thedevsaddam/renderer"
)

// api 관련 패키지

func MakeHandler() http.Handler {
	types.TodoMap = make(map[int]types.Todo)
	mux := mux.NewRouter()
	mux.Handle("/", http.FileServer(http.Dir("public")))
	mux.HandleFunc("/todos", GetTodoListHandler).Methods("GET")
	mux.HandleFunc("/todos", PostTodoHandler).Methods("POST")
	mux.HandleFunc("/todos/{id:[0-9]+}", RemoveTodoHandler).Methods("DELETE")
	mux.HandleFunc("/todos/{id:[0-9]+}", UpdateTodoHandler).Methods("PUT")

	return mux
}

func GetTodoListHandler(w http.ResponseWriter, r *http.Request) {
	list := make(types.Todos, 0)
	for _, todo := range types.TodoMap {
		list = append(list, todo)
	}

	sort.Sort(list)

	// json.NewEncoder(w).Encode(list)
	// w.WriteHeader(http.StatusOK)
	// 위 과정을 한 줄로 간편하게
	rnd := renderer.New()
	rnd.JSON(w, http.StatusOK, list)
}

func PostTodoHandler(w http.ResponseWriter, r *http.Request) {
	var todo types.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	types.LastID++
	todo.ID = types.LastID
	types.TodoMap[todo.ID] = todo
	rnd := renderer.New()
	rnd.JSON(w, http.StatusCreated, todo)
}

func RemoveTodoHandler(w http.ResponseWriter, r *http.Request) {
	// r에서 url variable 까서
	vars := mux.Vars(r)
	// id 체크하고
	id, err := strconv.Atoi(vars["id"])
	rnd := renderer.New()
	if err != nil {
		log.Fatal(err)
		rnd.JSON(w, http.StatusBadRequest, types.Success{Success: false})
		return
	}
	// map에서 id 조회후
	if _, ok := types.TodoMap[id]; ok {
		// 삭제한다.
		delete(types.TodoMap, id)
		rnd.JSON(w, http.StatusOK, types.Success{Success: true})
	} else {
		rnd.JSON(w, http.StatusOK, types.Success{Success: false})
	}

}

func UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	var newTodo types.Todo
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	rnd := renderer.New()
	if todo, ok := types.TodoMap[id]; ok {
		todo.Name = newTodo.Name
		todo.Completed = newTodo.Completed
		rnd.JSON(w, http.StatusOK, types.Success{Success: true})
	} else {
		rnd.JSON(w, http.StatusOK, types.Success{Success: false})
	}
}
