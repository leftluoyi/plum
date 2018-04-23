package services

//func AppendToBlockChain(block models.Block) {
//	dat, err := ioutil.ReadFile(os.Getenv("CHAIN_FILE"))
//	utils.Check(err)
//
//	var bc []models.Block
//	err = json.Unmarshal(dat, &bc)
//	utils.Check(err)
//
//	bc = append(bc, block)
//
//	result, err := json.Marshal(bc)
//	utils.Check(err)
//	err = ioutil.WriteFile("blockchain.json", result, 0644)
//	utils.Check(err)
//}
//
//func WriteBlockChain(blocks []models.Block) {
//	result, err := json.Marshal(blocks)
//	utils.Check(err)
//	err = ioutil.WriteFile("blockchain.json", result, 0644)
//}
//
//func GetBlockChain() []models.Block {
//	dat, err := ioutil.ReadFile(os.Getenv("CHAIN_FILE"))
//	utils.Check(err)
//
//	var bc []models.Block
//	err = json.Unmarshal(dat, &bc)
//	utils.Check(err)
//	return bc
//}
//
//func GetBlockChainString() string {
//	dat, err := ioutil.ReadFile(os.Getenv("CHAIN_FILE"))
//	utils.Check(err)
//
//	return string(dat)
//}
