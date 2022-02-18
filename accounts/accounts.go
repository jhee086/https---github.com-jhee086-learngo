package accounts

import (
	"errors"
	"fmt"
)

// Account struct
type Account struct {
	// 소문자로 시작하면 private이 되어 main에서 접근 불가능(unknown field)
	owner   string
	balance int
}

var errNoMoney = errors.New("Cant withdraw")

// constructor가 없기 때문에 func으로 construct 하거나 struct을 만듦
// 구글링: struct making function
// private한 struct만들고 public function 만들기
// export 하기 위해서는 func 첫글자 대문자

// NewAccount creates Account : make func that return objects
func NewAccount(owner string) *Account {
	// account 초기화
	account := Account{owner: owner, balance: 0}
	return &account // 실제 메모리 address를 return (return object)
	// 새로운 object return 하고 싶을 때 복사본 자체를 return
}

// method
/* receiver: struct의 첫 글자를 따서 소문자로 지어야 함
   (a Account): receiver(a의 type은 Account)
   func은 이름 전에 struct를 가지고 있다 */

// Deposit X amount on your account
// *Account : Go에게 account나 receiver를 복사하지 말고 실제 receiver를 주는 것
func (a *Account) Deposit(amount int) {
	fmt.Println("Gonna deposit", amount)
	a.balance += amount
}

// Balance of your account
func (a Account) Balance() int {
	return a.balance
}

// Withdraw X amount from your account
func (a *Account) Withdraw(amount int) error {
	/* error handling - Go에는 exception이 없음
	: error를 return 해주고 error 직접 체크 필요*/
	if a.balance < amount {
		return errNoMoney //errors.New("Cant withdraw you are poor")
	}
	a.balance -= amount
	return nil
}

// CangeOwner of the account
func (a *Account) ChangeOwner(newOwner string) {
	a.owner = newOwner
}

// owner 가져오기 : 변경필요없어 복사본 사용해도 됨
// Owner of the account
func (a Account) Owner() string {
	return a.owner
}

// String(): string으로 표현해서 출력 (자동으로 호출) -> Python의 __str__
// Go가 내부적으로 호출하는 method를 사용하는 방법
func (a Account) String() string {
	return fmt.Sprint(a.Owner(), "'s account.\nHas: ", a.Balance()) //"Whatever you want"
}

// dictionary - map
