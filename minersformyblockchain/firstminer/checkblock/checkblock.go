package checkblock

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	. "miner/firstminer/config"
	"miner/firstminer/createblock"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func calculateHash(block *createblock.Block) string {

	blockData := fmt.Sprintf("%d|%s|%d|%v|%d", block.Index, block.PrevHash, block.Timestamp, block.Transaction, block.Nonce)
	hashBytes := sha256.Sum256([]byte(blockData))
	blockHash := hex.EncodeToString(hashBytes[:])
	return blockHash
}

func CheckBlock(c *gin.Context) {
	block := &createblock.Block{}
	Index, err := strconv.Atoi(c.PostForm("Index"))
	if err != nil {
		fmt.Println("ошибка в index блока", Index)
		return
	}
	block.Index = Index
	Hash := c.PostForm("Hash")
	if Hash == "" {
		fmt.Println("ошибка в hash", Hash)
		return
	}
	block.Hash = Hash
	PrevHash := c.PostForm("PrevHash")
	block.PrevHash = PrevHash
	Timestamp, err := strconv.ParseInt(c.PostForm("Timestamp"), 10, 64)
	if err != nil {
		fmt.Println("ошибка в таймстампе", err)
		return
	}
	block.Timestamp = Timestamp
	Transaction := c.PostForm("Transaction")
	if Transaction == "" {
		fmt.Println("ошибка в Transaction", Transaction)
		return
	}
	block.Transaction = Transaction
	Nonce, err := strconv.Atoi(c.PostForm("Nonce"))
	if err != nil {
		fmt.Println("ошибка в index блока", Nonce)
		return
	}
	block.Nonce = Nonce

	hash := calculateHash(block)
	if hash == block.Hash {
		Add_block(block)
		c.JSON(http.StatusAccepted, gin.H{"status": "1"})

	} else {
		c.JSON(http.StatusAccepted, gin.H{"status": "0"})
	}
}

func Add_block(block *createblock.Block) {

	_, err := DB.Exec("INSERT INTO blocks(\"index\",hash, prev_hash, timestamp, \"transaction\", nonce) VALUES(?,?,?,?,?,?)", block.Index, block.Hash, block.PrevHash, block.Timestamp, block.Transaction, block.Nonce)
	if err != nil {
		fmt.Println("ошибка в добавлении блока здесь ", err)
	}
}
