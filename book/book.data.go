package book

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"
)

var bookMap = struct {
	sync.RWMutex
	m map[int]Book
}{m: make(map[int]Book)}

func init() {
	loadedBookMap, err := loadBookMap()
	if err != nil {
		log.Fatal(err)
		return
	}
	bookMap.m = loadedBookMap
}

func loadBookMap() (map[int]Book, error) {

	fileName := "books.json"
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file [%s] does not exist", fileName)
	}

	file, _ := ioutil.ReadFile(fileName)
	bookList := make([]Book, 0)
	err = json.Unmarshal([]byte(file), &bookList)
	if err != nil {
		log.Fatal(err)
	}

	bookListMap := make(map[int]Book)
	for i := 0; i < len(bookList); i++ {
		bookListMap[bookList[i].Id] = bookList[i]
	}

	return bookListMap, nil
}

func getBook(id int) *Book {
	bookMap.RLock()
	defer bookMap.RUnlock()
	if book, ok := bookMap.m[id]; ok {
		return &book
	}
	return nil
}

func removeBook(id int) {
	bookMap.Lock()
	defer bookMap.Unlock()
	delete(bookMap.m, id)
}

func getBookList() []Book {
	bookMap.RLock()

	books := make([]Book, 0, len(bookMap.m))
	for _, value := range bookMap.m {
		books = append(books, value)
	}
	bookMap.Unlock()
	return books
}

func getBookIds() []int {
	bookMap.RLock()
	bookIds := []int{}
	for key := range bookMap.m {
		bookIds = append(bookIds, key)
	}
	bookMap.RUnlock()
	sort.Ints(bookIds)
	return bookIds
}

func getNextBookID() int {
	bookIDs := getBookIds()
	return bookIDs[len(bookIDs)-1] + 1
}

func updateBook(book Book) (int, error) {
	id := -1
	if book.Id > 0 {
		dbProduct := getBook(book.Id)

		if dbProduct == nil {
			return 0, fmt.Errorf("book id [%d] doesn't exist", book.Id)
		}
		id = book.Id
	} else {
		id = getNextBookID()
		book.Id = id
	}
	bookMap.Lock()
	bookMap.m[id] = book
	bookMap.Unlock()
	return id, nil
}
