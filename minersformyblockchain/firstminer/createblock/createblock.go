package createblock

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	. "miner/firstminer/config"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Block struct {
	Index       int
	Hash        string
	PrevHash    string
	Timestamp   int64
	Transaction string
	Nonce       int
}

// вычисление nonce для хэша
func ProofOfWork(block *Block, difficulty int) (int, string) {
	nonce := 0
	for {
		block.Nonce = nonce
		fmt.Println(nonce)
		hash := calculateHash(block)
		if hasLeadingZeros(hash, difficulty) {
			return nonce, hash
		}
		nonce++
	}
}

func calculateHash(block *Block) string {

	blockData := fmt.Sprintf("%d|%s|%d|%v|%d", block.Index, block.PrevHash, block.Timestamp, block.Transaction, block.Nonce)
	hashBytes := sha256.Sum256([]byte(blockData))
	blockHash := hex.EncodeToString(hashBytes[:])
	fmt.Println(blockHash)
	return blockHash
}

func hasLeadingZeros(hash string, difficulty int) bool {
	prefix := ""
	for i := 0; i < difficulty; i++ {
		prefix += "0"
	}
	return len(hash) >= difficulty && hash[:difficulty] == prefix
}

// AdDBlock добавляет новый блок, включающий указанные транзакции, в блокчейн.
func CreateBlock(c *gin.Context, name string) {
	fmt.Println("дошел")
	txID := c.PostForm("txID")
	if txID == "" {
		fmt.Println("Пустой txID")
		return
	}
	difficulty := 2

	row := DB.QueryRow("SELECT \"index\", hash FROM blocks ORDER BY \"index\" DESC LIMIT 1")
	var lastIndex int
	var lastHash string
	err := row.Scan(&lastIndex, &lastHash)
	newIndex := 0
	prevHash := ""
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Ошибка в запросе SELECT")
			return
		}
	} else {
		newIndex = lastIndex + 1
		prevHash = lastHash
	}
	timestamp := time.Now().Unix()
	block := &Block{
		Index:       newIndex,
		PrevHash:    prevHash,
		Timestamp:   timestamp,
		Transaction: txID,
	}
	Nonce, Hash := ProofOfWork(block, difficulty)
	block.Nonce = Nonce
	block.Hash = Hash

	fmt.Println("отправляю майнерам")
	err = PostMiners(block, name)
	if err != nil {
		fmt.Println("Ошибка в функции PostMiners ", err)
		return
	}

}

func PostMiners(block *Block, name string) error {
	type minerstruct struct {
		ip   string
		port int
	}

	rows, err := DB.Query("SELECT ip,port FROM miners")
	if err != nil {
		return err
	}
	fmt.Println(block)
	for rows.Next() {
		miner := minerstruct{}
		err := rows.Scan(&miner.ip, &miner.port)
		if err != nil {
			fmt.Println("Проблема", err)
		}
		form := url.Values{}
		form.Add("Index", strconv.Itoa(block.Index))
		form.Add("Hash", block.Hash)
		form.Add("PrevHash", block.PrevHash)
		form.Add("Timestamp", strconv.FormatInt(block.Timestamp, 10))
		form.Add("Transaction", block.Transaction)
		form.Add("Nonce", strconv.Itoa(block.Nonce))
		form.Add("Name", name)
		address := fmt.Sprintf("http://%s:%d/createdblock", miner.ip, miner.port)
		fmt.Println("Отправляем на:", address)

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.PostForm(address, form)
		if err != nil {
			fmt.Println("ошибка отправки на", address, ":", err)
			return err
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
	return nil
}
