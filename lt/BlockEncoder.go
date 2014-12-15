package lt

import (
	"fmt"
	"io"
	"math"
)

// BlockEncoder 实现LT编码的块编码器
type BlockEncoder struct {
	fileSize  uint32  // 原始文件大小(字节)
	blockSize uint32  // 编码块总大小(数据+头部)
	dataSize  uint32  // 纯数据块大小

	fileData []byte     // 文件数据缓存
	nBlocks  uint32     // 总数据块数
	planner  BlockPlanner // 块规划器
}

// uncodedBlock 表示未编码的原始数据块
type uncodedBlock []byte

// NewEncoder 创建新的LT编码器
// buf: 输入数据流
// size: 输入数据大小
// datablockSize: 每个数据块的大小
// initSeed: 随机种子
func NewEncoder(buf io.Reader, size uint64, datablockSize uint32, initSeed uint32) *BlockEncoder {
	// 读取输入流数据
	fileBuf := make([]byte, size)
	var bufPtr uint64
	for bufPtr < size {
		nRead, err := buf.Read(fileBuf[bufPtr:])
		if err != nil {
			fmt.Errorf(err.Error())
			return nil
		}
		bufPtr += uint64(nRead)
	}

	// 数据填充对齐到块大小
	padNum := size % uint64(datablockSize)
	if padNum != 0 {
		padNum = uint64(datablockSize) - padNum
		pads := make([]byte, padNum)
		fileBuf = append(fileBuf, pads...)
	}

	// 计算总块数
	nb := uint32(math.Ceil(float64(size)/float64(datablockSize)))
	return &BlockEncoder{
		blockSize: datablockSize + BLOCK_HEADER_SIZE,
		dataSize:  datablockSize,
		fileSize:  uint32(size),

		fileData: fileBuf,
		nBlocks:  nb,
		planner:  NewBlockPlanner(nb, initSeed),
	}
}

// getBlock 获取指定索引的原始数据块
func (enc *BlockEncoder) getBlock(bnum uint32) uncodedBlock {
	startIdx := bnum * enc.dataSize
	return enc.fileData[startIdx : startIdx+enc.dataSize]
}

// NextCodedBlock 生成下一个编码块
func (enc *BlockEncoder) NextCodedBlock() CodedBlock {
	// 获取需要组合的块列表和当前种子
	blockList, currSeed := enc.planner.NextBlockList()
	// 创建空块用于累积
	accum_block := uncodedBlock(make([]byte, enc.dataSize))

	// 对选中的块进行异或组合
	for _, blockIdx := range blockList {
		accum_block.xorBlock(enc.getBlock(blockIdx))
	}

	// 返回编码块
	ans := CodedBlock{
		fileSize:  enc.fileSize,
		blockSize: enc.blockSize,
		seed:      currSeed,
		data:      accum_block}
	return ans
}

// xorBlock 对两个块进行异或操作(原地修改x)
func (x uncodedBlock) xorBlock(y uncodedBlock) {
	if len(x) != len(y) {
		panic("xoring unequal length lists")
	}
	for i, yb := range y {
		x[i] ^= yb
	}
}
	

