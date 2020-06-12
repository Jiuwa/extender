# extender

#### `main.go`

init()以当前时间作种子获取随机数
main()给httprouter设定了controller的Index,Filter,Prioritize()，输出开始提示以及监听8888端口（防止端口冲突）

#### `predicates.go`

filter()：对该pod，遍历所有节点，使用podFitsOnNode来判断pod是否node，匹配的与不匹配的分开存储。
podFitsOnNode()：对一个node和一个pod使用PodNameFitPredicate()来判断是否匹配，匹配结果整理之后返回给filter。
PodNameFitPredicate()：匹配判断的核心逻辑r若pod的名称长度大于30，则这个节点可以被批准，否则，不被批准。

#### `priorities.go`

prioritize()：对已过滤的节点进行打分。打分的规则是被调度的pod名称长度加上0~最大优先级的随机数，作为该节点的分数。

#### `routers.go`

Index,Filter,Prioritize()作为接口函数，将过滤和打分结果处理后按照kubernetes需要的格式返回

