# 日志定义
## 日志等级
	1. Panic 致命信息，导致业务不可用
	2. Error 错误信息， 比较严重的信息，可能引起业务的批量失败，
	3. Warning 警告信息，普通的错误信息，对业务影响较小
	4. Info  普通信息，用于记录一些必要信息，可以用来分析系统的数据，一般生产环境的最低级别
	5. Debug 调试信息，用于记录内部的流程，数据，便于定位问题，一般在生产环境会被关闭