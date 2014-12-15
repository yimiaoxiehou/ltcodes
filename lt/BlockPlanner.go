package lt

// BlockPlanner 负责规划LT编码中使用的数据块组合
type BlockPlanner struct {
	nblocks uint32  // 总数据块数
	rando   RandGen // 随机数生成器
	sol     Soliton // Soliton分布生成器
}

// NewBlockPlanner 创建新的块规划器
// nb: 总数据块数
// seed: 随机种子
func NewBlockPlanner(nb uint32, seed uint32) BlockPlanner {
	return BlockPlanner{
		nblocks: nb,
		rando:   &LinearGen{seed: seed},
		sol:     NewSoliton(nb),
	}
}

// CurrSeed 获取当前随机种子
func (planner BlockPlanner) CurrSeed() uint32 {
	return planner.rando.getSeed()
}

// NextBlockList 生成下一个编码块列表
// 返回: 块索引列表和当前使用的随机种子
func (planner BlockPlanner) NextBlockList() (blockList []uint32, currSeed uint32) {
	currSeed = planner.rando.getSeed() // 获取当前种子
	nToCode := planner.sol.generate(planner.rando) // 根据Soliton分布生成要编码的块数
	blockList = make([]uint32, nToCode) // 初始化块列表
	var n uint // 已添加块计数器

addBlock:
	for n < nToCode {
		nextBlock := planner.rando.nextInt() % planner.nblocks // 随机选择块
		// 检查是否已选择过该块
		for _, currBlock := range blockList[:n] {
			if currBlock == nextBlock {
				continue addBlock // 跳过重复块
			}
		}
		blockList[n] = nextBlock // 添加新块
		n++
	}
	return
}
