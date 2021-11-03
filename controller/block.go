package controller

import (
	"encoding/json"
	"net/http"

	"github.com/eduardohoraciosanto/simple-blockchain/pkg/blockchain"
	"github.com/eduardohoraciosanto/simple-blockchain/viewmodels"
)

type BlockController struct {
	Service blockchain.Service
}

//Health is the handler for the health endpoint
func (c *BlockController) GetBlockchain(w http.ResponseWriter, r *http.Request) {

	//using lower level pkg to do the logic
	blockchain := c.Service.GetBlockchain()

	viewmodels.RespondWithData(w, http.StatusOK, blockchain)
}

func (c *BlockController) InsertDataBlock(w http.ResponseWriter, r *http.Request) {
	vm := viewmodels.InsertDataRequest{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&vm)
	if err != nil {
		viewmodels.RespondBadRequest(w)
		return
	}
	//using lower level pkg to do the logic
	if err := c.Service.InsertDataBlock(blockchain.BlockData{
		BMP: vm.BlockData.BMP,
	}); err != nil {
		viewmodels.RespondInternalServerError(w)
		return
	}

	viewmodels.RespondWithData(w, http.StatusAccepted, nil)
}
