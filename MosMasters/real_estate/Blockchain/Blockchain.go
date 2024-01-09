package Blockchain

import (
	"log"

	"github.com/boltdb/bolt"
	"github.com/VladislavSCV/Blockchain_/MosMasters/real_estate/Blockchain/Block"
)

const dbFile = "blockchain.db"
const blocksBucket = ""

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (i *BlockchainIterator) Next() *Block.Block {
	var block *Block.Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = Block.DeserializeBlock(encodedBlock)

		return nil
	})
	if err != nil {
		log.Panic("BlockchainIterator error:", err)
	}

	i.currentHash = block.PrevBlockHash

	return block
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}

func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})
	if err != nil{
		log.Panic("Ошибка в добавлении блока")
	}

	newBlock := Block.NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		err = b.Put([]byte("l"), newBlock.Hash)
		bc.tip = newBlock.Hash
		if err != nil {
			log.Panic("Ошибка в destruct")
		}

		return nil
	})
}

func NewGenesisBlock() *Block.Block {
    return Block.NewBlock("Genesis Block", []byte{})
}

func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			err = b.Put(genesis.Hash, genesis.Serialize())
			err = b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
			if err != nil {
				log.Panic("Ошибка в создании блокчейна")
			}
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	if err != nil {
		log.Panic("Ошибка в создании блокчейна")
	}

	bc := Blockchain{tip, db}

	return &bc
}

