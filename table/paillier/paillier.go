package paillier

import (
	"fmt"
	"math/big"
	"math/rand"
	"time"
)

func IsPrime(n int) bool {
	if n == 1 {
		return false
	}
	for i := 2; i < n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func prime(start int) int {
	for num := start; num <= start+200000; num++ {
		if IsPrime(num) {
			return num
		} else {
			continue
		}
	}
	return -1
}

func genPrime() (int64, int64) {
	rand.Seed(time.Now().Unix())
	a, b := prime(rand.Intn(200)+100), prime(rand.Intn(400)+200)
	c := int64(a)
	d := int64(b)
	return c, d
}

func gcd(x, y *big.Int) *big.Int {
	tmp := big.NewInt(1).Mod(x, y)
	flag := tmp.Cmp(big.NewInt(0))
	if flag > 0 {
		return gcd(y, tmp)
	} else {
		return y
	}
}

func lcm(x, y *big.Int) *big.Int {
	tmp := big.NewInt(1).Mul(x, y)
	return big.NewInt(1).Div(tmp, gcd(x, y))
}

type KeyPaillier struct {
	p       *big.Int
	q       *big.Int
	lamnda  *big.Int
	n       *big.Int
	nsquare *big.Int
	g       *big.Int
}

func KeyGenPaillier() KeyPaillier {
	var paillier KeyPaillier
	a, b := genPrime()
	p := big.NewInt(a)
	q := big.NewInt(b)
	n := big.NewInt(1).Mul(p, q)
	paillier.p = p
	paillier.q = q
	paillier.n = n
	nsquare := big.NewInt(1).Mul(n, n)
	paillier.nsquare = nsquare
	g := big.NewInt(2)
	paillier.g = g
	p_1 := big.NewInt(1).Add(p, big.NewInt(-1))
	q_1 := big.NewInt(1).Add(q, big.NewInt(-1))
	lamnda := lcm(p_1, q_1)
	paillier.lamnda = lamnda
	return paillier
}

func Print(paillier KeyPaillier) string {
	fmt.Print("p = ")
	fmt.Println(paillier.p)
	fmt.Print("q = ")
	fmt.Println(paillier.q)
	fmt.Print("lamnda = ")
	fmt.Println(paillier.lamnda)
	fmt.Print("n = ")
	fmt.Println(paillier.n)
	fmt.Print("nsquare = ")
	fmt.Println(paillier.nsquare)
	fmt.Print("g = ")
	fmt.Println(paillier.g)
	return "{p = " + paillier.p.String() + ", q = " + paillier.q.String() + ", lamnda = " + paillier.lamnda.String() + ", n = " + paillier.n.String() + ", nsquare = " + paillier.nsquare.String() + ", g = " + paillier.g.String() + "}"
}

func Encryption(paillier KeyPaillier, m *big.Int) *big.Int {
	rand.Seed(time.Now().Unix())
	rint := int64(rand.Intn(20000))
	r := big.NewInt(rint)
	//     fmt.Print("r = ")
	//     fmt.Println(r)

	temp1 := big.NewInt(1).Exp(paillier.g, m, nil)
	temp2 := big.NewInt(1).Mod(temp1, paillier.nsquare)
	temp3 := big.NewInt(1).Exp(r, paillier.n, nil)
	temp4 := big.NewInt(1).Mod(temp3, paillier.nsquare)
	temp5 := big.NewInt(1).Mul(temp2, temp4)
	return big.NewInt(1).Mod(temp5, paillier.nsquare)
}

func Decryption(paillier KeyPaillier, c *big.Int) *big.Int {
	temp1 := big.NewInt(1).Exp(c, paillier.lamnda, nil)
	temp2 := big.NewInt(1).Mod(temp1, paillier.nsquare)
	temp3 := big.NewInt(1).Add(temp2, big.NewInt(-1))
	temp4 := big.NewInt(1).Div(temp3, paillier.n)

	temp5 := big.NewInt(1).Exp(paillier.g, paillier.lamnda, nil)
	temp6 := big.NewInt(1).Mod(temp5, paillier.nsquare)
	temp7 := big.NewInt(1).Add(temp6, big.NewInt(-1))
	temp8 := big.NewInt(1).Div(temp7, paillier.n)
	temp9 := big.NewInt(1).ModInverse(temp8, paillier.n)
	temp10 := big.NewInt(1).Mul(temp4, temp9)
	temp11 := big.NewInt(1).Mod(temp10, paillier.n)
	return temp11
}

func addMiWenPaillier(paillier KeyPaillier, c1 *big.Int, c2 *big.Int) *big.Int {
	result := big.NewInt(1).Mul(c1, c2)
	return Decryption(paillier, result)
}

func GenKey(keys []int64) (KeyPaillier, error) {
	var paillier KeyPaillier
	if len(keys) != 6 {
		return paillier, fmt.Errorf("Incorrect number of parameters for KEY")
	}
	paillier.p = big.NewInt(keys[0])
	paillier.q = big.NewInt(keys[1])
	paillier.lamnda = big.NewInt(keys[2])
	paillier.n = big.NewInt(keys[3])
	paillier.nsquare = big.NewInt(keys[4])
	paillier.g = big.NewInt(keys[5])
	return paillier, nil
}
