# ASM磁盘管理
## 相关概念
    磁盘组， 磁盘，故障组， 分配单元， ASM文件，I/O分布，Rebalance， ASM磁盘组的管理
### 1.ASM 磁盘组
    一个ASM磁盘组由过多个ASM磁盘组成
    一个磁盘组内可以存放多个数据文件，
    一个数据文件仅仅只能位于一个磁盘组内，不能跨磁盘组
    多个数据库可以共享相同的或多个磁盘组
    磁盘组的冗余类型可以分为三类：标准冗余，高度冗余，外部冗余
    对于已创建的磁盘组，不能够更改其冗余级别，如要更改，需要删除该磁盘组后再重新创建

### 2.ASM 磁盘
    ASM磁盘通过标准的OS接口来访问，由Oracle用户来读写，在聚集的所有节点可以被访问
    ASM磁盘在不同的节点可以使用不同的名字
    ASM磁盘可以使网络文件系统
    ASM磁盘上的对象被冗余保护
    每一个ASM磁盘的第1块用于定义磁盘的头部信息，ASM磁盘名字编号，创建的时间戳等
    ASM文件会均匀分布在一个ASM组内的各个磁盘中

### 3.ASM 故障组
    一个磁盘组可以由两个或多个故障组组成
    一个故障组由一个或多个ASM磁盘组成
    故障组提供了共享相同资源的冗余，我们可以这样来理解标准冗余
    假定有磁盘组DG1,且创建了两个故障组fgroup1,fgroup2,每个故障组由2个ASM磁盘组成,则对标准冗余而言,两个故障组互为镜像
    failgroup1 --> asmdiskA , asmdiskB
    failgroup2 --> asmdiskC , asmdiskD
    假定文件datafileA大小为4MB,则4个extent均匀分布到asmdiskA,asmdiskB,同样asmdiskC,asmdiskD也包含该文件的1至4个extent
    即只要有一个extent在故障组fgroup1中存在，必定有一个镜像的extent存在于fgroup2中，反之亦然，两个extent互为镜像。
    当一个故障组中的某个磁盘损坏，假定为asmdiskA ，则asmdiskA中原来保存的extent将会从failgroup2中复制到asmdiskB中。
    
    总之，故障组failgroup1和failgroup2必定有相同的extent副本
    标准冗余至少需要2个故障组，
    高度冗余则至少需要3个故障组。
    事实上对于未明确指定故障组的情况下，一个标准冗余至少需要2个asm磁盘，而高度冗余至少需要3个asm磁盘

### 4.分配单元
    ASM磁盘的最小粒度是分配单元，大小默认是1M，也可设置为128K进行细粒度访问
    支持粗粒度和细粒度分配单元进行读写来实现装载平衡和减少延迟
    ASM文件由一些分配单元的集合组成

### 5.ASM 文件 
    对Oracle自身而言，实际上与标准的文件并没有太多区别
    ASM文件一般位于磁盘组内创建的子目录内，磁盘组以加号开头，相当于Linux系统的根目录
    如+DG1/oradb/datafile/system.258.346542
    ASM可以为控制文件，数据文件，联机日志文件，参数文件，归档日志，备份等
    不支持trace文件，可执行文件，OCR，Votingdisk等，注：Oracle 11g R2可支持
    使用extent maps来记录文件到磁盘的映射

### 6.I/O分布
    可以使用条带化和镜像来保护数据
    文件被平均分布在一个组内的所有磁盘中
    磁盘的添加与删除，ASM会自动重新分配AU，因此也不存在碎片的问题
    将I/O分配到不同的磁盘控制器提高了读写数据

### 7.Rebalance
    ASM 文件被均衡地分布在一个磁盘组的所有磁盘中
    磁盘添加时，当前磁盘组加载的所有磁盘中共享的部分extent将会被移植到新的磁盘中,直到重新分布完成才正常提供I/O均衡
    磁盘删除或故障时，删除磁盘或故障磁盘的extent将会被均匀的分布到剩余的磁盘中
    未使用force关键字drop磁盘操作,该磁盘上所有数据rebalance完毕后才被释放.即完毕后磁盘脱机,置磁盘头部状态为former
    总之,任意存储性质改变(磁盘增加,删除,故障)都将导致rebalance,且由asm自动完成,无需人工干预,在一个时间段通常会锁定一个盘区

### 8.ASM磁盘组的管理
    通常建议创建两个磁盘组，一个用于保存数据文件，一个用于保存闪回，备份恢复使用
    Flash Recovery Area 的大小取决于闪回内容需要保留的时间长短
    尽可能将数据区与闪回区使用不同的物理通道
    尽可能一次性mount所有需要用到的磁盘
    建议使用性能，磁盘大小相近的磁盘。假定两个故障组FG1，FG2各使用一块磁盘，则FG1内的磁盘应保持与FG2内的磁盘大小相同，
    否则会以最小的磁盘空间作为可使用空间
    
### 9.ASM磁盘组管理原则
    添加或删除磁盘的影响
    当发生添加/删除磁盘组中磁盘的操作时，ASM能够自动平衡。
    对于普通的删除操作(无force选项)，被删除的磁盘在该上数据被有效处理前并不会立刻释放
    新增磁盘时，在重分配工作完成前，该盘也不会承担I/O负载的工作

### 10.ASM如何处理磁盘故障
    ASM磁盘组大致有二：普通组和故障组，后者与ASM的冗余方式有所关联。
    普通磁盘组就是标准的存储单元，ASM可以向其可访问的磁盘组中读写数据，failure磁盘组是为了提高数据的高可用性。
    ASM中的磁盘冗余策略非常简单，概要成三类：外部冗余、标准冗余和高度冗余。其中，外部冗余和failure组无关。
    如果设置了标准冗余或者高度冗余，那么该磁盘组就必须要有故障组。 
    对于标准冗余，ASM要求该磁盘组至少要拥有两个failure磁盘组，即提供双倍镜像保护，对于同一份数据，将有主从两份镜像。 
    并且ASM通过算法来自动确保主、从镜像不会存在于同一份failure磁盘组，这样就保障了就算整个failure磁盘组都损坏，数据也不会丢失。
    ASM中镜像单位不是磁盘，也不是块，而是一种AU的单位，该单位大小默认是1M。
    至于高度冗余，它至少需要三个failure磁盘组，也就是一份AU有一主多从的镜像，理论上将更加安全。
    外部冗余的话磁盘属于磁盘组，内部冗余的话，磁盘属于磁盘组的同时，还属于而且仅属于某个failure磁盘组。
    如果磁盘发生损坏，那么损坏的磁盘默认自动offlice并被drop掉，不过该磁盘所在的磁盘组仍将保持MOUNT状态。 
    如果该组有镜像的话，那么应用不会有影响，镜像盘将自动实现接管--只要不是所有failure磁盘组都损坏掉，否则的话，该磁盘组将自动DISMOUNT 
    举个例子吧，某标准冗余的failure组有6个盘(对应6个裸设备)，假如说此时坏了一块盘，没关系，操作继续，坏了那块会被自动dropped，剩下的5块盘仍然能够负担起正常的读写操作。

    
## 相关命令
    select instance_name from v$instance;
    
    查看当前已经存在的ASM磁盘组
    set pagesize 1000 linesize 500
    select GROUP_NUMBER,NAME,TOTAL_MB,FREE_MB from v$asm_diskgroup;
    
    查看ASM磁盘的冗余策略
    select state,name,type from v$asm_diskgroup;
    
    查看磁盘信息
    select GROUP_NUMBER,DISK_NUMBER,TOTAL_MB,FREE_MB,NAME from v$asm_disk;
    
    添加磁盘
    alter diskgroup DATA add disk 'ORCL:DATA1' rebalance power 10;

    监控asm磁盘组平衡速度
    select * from v$asm_operation;   EST_MINUTES=0的时候 表示完成

    备份磁盘头信息
    su - oracle
    mkdir -p /u01/app/oracle/diskheader
    cd /u01/app/oracle/diskheader
    [oracle@testya diskheader]$ kfed read /dev/oracleasm/disks/DATA > DATAheader.txt
    [oracle@testya diskheader]$ kfed read /dev/oracleasm/disks/DATA1 > DATA1header.txt

    ---删除磁盘组
    select GROUP_NUMBER,DISK_NUMBER,TOTAL_MB,FREE_MB,NAME from v$asm_disk; 

    
    