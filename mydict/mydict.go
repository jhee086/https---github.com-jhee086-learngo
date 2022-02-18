package mydict

import "errors"

// Dictionary type -> map[string]string에 대한 alias(가명)를 만듦
type Dictionary map[string]string

var (
	errNotFound   = errors.New("Not Found")
	errWordExits  = errors.New("That word already exists")
	errCantUpdate = errors.New("Cant update non-existing word")
)

// method를 가질 수 있다.
func (d Dictionary) Search(word string) (string, error) {
	// key의 존재여부를 알려주는 방법 있음
	// map의 key를 호출하면 2개의 값을 얻음 string boolean
	value, exists := d[word]
	if exists {
		return value, nil
	}
	return "", errNotFound
}

// Add a word to the dictionary
func (d Dictionary) Add(word, def string) error {
	// 추가하려는 단어가 있으면 에러! 없으면 추가
	_, err := d.Search(word) // 에러가 없으면 단어가 존재한다는 뜻
	switch err {
	case errNotFound:
		d[word] = def
	case nil:
		return errWordExits
	}
	return nil
	// if err == errNotFound {
	// 	d[word] = def
	// } else if err == nil {
	// 	return errWordExits
	// }
	// return nil
}

// Update a word
func (d Dictionary) Update(word, definition string) error {
	_, err := d.Search(word)
	switch err {
	case nil:
		d[word] = definition
	case errNotFound:
		return errCantUpdate
	}
	return nil
}

// Delete a word
func (d Dictionary) Delete(word string) {
	// delete 함수는 아무것도 return하지 않고
	// 특정한 key가 없으면 아무것도 하지 않을 것
	delete(d, word)
}
