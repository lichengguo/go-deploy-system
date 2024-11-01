<template>
    <div>
        <div class="crumbs">
            <el-breadcrumb separator="/">
                <el-breadcrumb-item>
                    <i class="el-icon-lx-calendar"></i> 发布项目
                </el-breadcrumb-item>
            </el-breadcrumb>
        </div>

        <el-row>
            <!-- 项目信息 -->
            <el-col :span="11">
                <el-card shadow="hover" style="margin-bottom: 20px" >
                    <div style="margin-bottom:15px;font-size:24px;color: blue;">
                        <p>项目信息</p>
                        <hr>
                    </div>
                    <div>
                        <div style="margin-bottom:15px">
                            <p style="margin-bottom:5px;margin-top:20px">选择项目</p>
                            <el-select v-model="values" placeholder="请选择发布项目" @change="selectData()">
                                <el-option v-for="item in releasesList" :key="item.id" :label="item.deploy_name" :value="item.id" >
                                    <span style="float:left; font-size:13px">{{ item.deploy_name }}</span>
                                </el-option>
                            </el-select>

                            <p style="margin-bottom:5px;margin-top:20px">备注(默认为Git提交信息)</p>
                            <el-input style="margin-bottom:5px" v-model="getData.GitInfo" placeholder="备注信息(暂无内容)"></el-input>
                            
                            <p style="margin-bottom:5px;margin-top:20px">Git相关信息</p>
                            <el-input style="margin-bottom:5px" :disabled="true" v-model="getData.GitUrl" placeholder="GitUrl(暂无内容)"></el-input>
                            <el-input style="margin-bottom:5px" :disabled="true" v-model="getData.GitHash" placeholder="Git指针(暂无内容)"></el-input>
                            <el-input style="margin-bottom:5px" :disabled="true" v-model="getData.GitEmail" placeholder="提交邮箱(暂无内容)"></el-input>

                            <p style="margin-bottom:5px;margin-top:20px">目标服务器:发布目录</p>
                            <el-input style="margin-bottom:5px" :disabled="true" v-model="getData.ServeripPath" placeholder="暂无内容"></el-input>
                        </div>
                    </div>
                </el-card>
            </el-col>

            <!-- 发布文件 -->
            <el-col :span="11" :offset="2">
                <el-card shadow="hover" class="mgb20" >
                    <div style="margin-bottom:15px;font-size:24px;color: red;">
                        <p>发布文件</p>
                        <hr>
                    </div>
                    <div>
                        <el-input type="textarea" placeholder="* 代表发布整个目录(请谨慎使用), 多个文件用换行隔开" :rows="17" v-model="getData.GitFileList"></el-input>
                        <el-button type="danger" style="width:80px;margin-top:20px" @click="deployCode">发布代码</el-button>
                    </div>
                </el-card>
            </el-col>

            <!-- 日志展示 -->
            <hr>
            <hr>
            <hr>
            <el-col :span="24">
                <div style="margin-bottom:10px;font-size:18px;color: blue;margin-top:20px">
                    <p>最近3次发布日志</p>
                </div>
                <el-table :data="logList" border class="table" ref="multipleTable" header-cell-class-name="table-header">
                    <el-table-column prop="ID" label="ID"></el-table-column>
                    <el-table-column prop="deployment_id" label="项目ID" align="center"></el-table-column>
                    <el-table-column prop="deployment_name" label="项目名称" align="center"></el-table-column>
                    <el-table-column prop="deployment_user_name" label="发布用户" align="center"></el-table-column>
                    <el-table-column prop="deployment_commit" label="发布信息备注" align="center"></el-table-column>
                    <el-table-column prop="deployment_status" label="发布状态" align="center">
                        <template slot-scope="scope">
                            <span v-if="scope.row.deployment_status==1" style="color:blue">成功</span>
                            <span v-if="scope.row.deployment_status==2" style="color:red">失败</span>
                        </template>
                    </el-table-column>
                    <el-table-column prop="UpdatedAt" label="发布时间" align="center">
                        <template slot-scope="scope">
                            {{dateFormat(scope.row.UpdatedAt)}}
                        </template>
                    </el-table-column>
                    <el-table-column label="操作" width="210" align="center">
                        <template slot-scope="scope">
                            <el-button type="danger" @click="rollbackinfo(scope.row.ID)">回滚</el-button>
                            <el-button type="primary"  @click="details(scope.row)">详情</el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </el-col>

            <!-- 日志详情弹框 -->
            <el-dialog title="日志详情"  :visible.sync="dialogFormVisible" top="0" center>
                <el-form :model="form" size="mini" :disabled=true label-width="100px">
                    <el-form-item label="ID">
                        <el-input v-model="form.ID"  ></el-input>
                    </el-form-item>
                    <el-form-item label="项目ID">
                        <el-input v-model="form.deployment_id"></el-input>
                    </el-form-item>
                    <el-form-item label="项目名称">
                        <el-input v-model="form.deployment_name"></el-input>
                    </el-form-item>
                    <el-form-item label="发布用户">
                        <el-input v-model="form.deployment_user_name"></el-input>
                    </el-form-item>
                    <el-form-item label="发布信息备注">
                        <el-input v-model="form.deployment_commit" ></el-input>
                    </el-form-item>
                    <el-form-item label="发布文件">
                        <el-input v-model="form.deployment_file_list" type="textarea" ></el-input>
                    </el-form-item>
                    <el-form-item label="发布状态">
                        <el-input v-model="form.deployment_status"></el-input>
                    </el-form-item>
                    <el-form-item label="失败详情">
                        <el-input v-model="form.deployment_fail_info"></el-input>
                    </el-form-item>
                    <el-form-item label="发布时间">
                        <el-input v-model="form.UpdatedAt"></el-input>
                    </el-form-item>
                    <el-form-item label="回滚指针">
                        <el-input v-model="form.git_head"></el-input>
                    </el-form-item>
                </el-form>
            </el-dialog>
        </el-row>
    </div>
</template>

<script>
import config from '../common/config.vue';
let loading // 发布动画变量

export default {    
    data(){
        return {
            dialogFormVisible: false,   // 日志弹框标志位
            form: {},                   // 日志弹框数据
            logList:[],                 // 日志列表数据
            values:"",                  // 选择的项目ID
            releasesList:[],            // 选择发布项目列表信息
            IPList:[],                  // 服务器IP列表数据
            // 项目信息数据
            getData:{
                GitInfo:'',
                GitUrl:'',
                GitHash:'',
                GitEmail:'',
                GitFileList:'', 
                ServerPath:'',
                ServeripPath:'',
            },
            // 发送给服务端的数据
            sendData:{
                deployment_id:0,
                deployment_file_list:'',
                deployment_commit:'',
                git_head:'',
            }
        }
    },

    methods:{
        // 格式化时间显示函数
        dateFormat(datetime){
            var date = new Date(datetime*1000); //时间戳为10位需*1000，时间戳为13位的话不需乘1000
            var year = date.getFullYear(),
                month = ("0" + (date.getMonth() + 1)).slice(-2),
                sdate = ("0" + date.getDate()).slice(-2),
                hour = ("0" + date.getHours()).slice(-2),
                minute = ("0" + date.getMinutes()).slice(-2),
                second = ("0" + date.getSeconds()).slice(-2);
            var result = year + "-"+ month +"-"+ sdate +" "+ hour +":"+ minute +":" + second; // 拼接
            return result; // 返回
        },

        // 日志详情
        details(row){
            this.dialogFormVisible = true
            this.form.ID = row.ID
            this.form.deployment_id = row.deployment_id
            this.form.deployment_name = row.deployment_name
            this.form.deployment_user_name = row.deployment_user_name
            this.form.deployment_file_list = row.deployment_file_list
            this.form.deployment_commit = row.deployment_commit
            if (row.deployment_status == 1) {
                this.form.deployment_status = '成功'
            } else {
                this.form.deployment_status = '失败'
            }
            this.form.deployment_fail_info = row.deployment_fail_info
            this.form.git_head = row.git_head
            this.form.UpdatedAt = this.dateFormat(row.UpdatedAt)
        },

        // 获取日志列表 只显示最后3条
        getLogList(){
            this.$axios.get(config.urlpath+'/deploymentlogs?pagesize=3&page=1').then(response => {
                if (response.data.status != 200) {
                    this.$message({message:'获取最近3条日志失败,请联系管理员', showClose: true, duration: 0,type: 'error'});
                }
                this.logList = response.data.data
            }).catch(error => {
                alert('获取最近3条日志失败,请联系管理员')
            })
        },

        // 获取项目列表信息
        getreleasesList(){
            this.$axios.get(config.urlpath+'/releases').then(response => {
                if(response.data.status == 200){
                    this.releasesList = response.data.data;
                }else{
                    this.$message({
                        type: 'error',
                        message: response.data.message
                    })
                }
                return
            }).catch(function (error) {
                this.$message({
                    type: 'error',
                    message: error
                })
                return
            });
        },

        // 选择具体项目，获取相关信息
        selectData(){
            if(this.values>0){
                this.$set(this.getData, "GitInfo",  "正在拉取中...");
                this.$set(this.getData, "GitUrl",  "正在拉取中...");
                this.$set(this.getData, "GitHash",  "正在拉取中...");
                this.$set(this.getData, "GitEmail",  "正在拉取中...");
                this.$set(this.getData, "ServeripPath",  "正在拉取中...");
                
                this.$axios.put(config.urlpath+'/release/gitpull/'+this.values).then(response =>{
                    if(response.data.status == 200){
                        this.getData.GitInfo = response.data.data.GitInfo
                        this.getData.GitUrl =  response.data.data.GitUrl
                        this.getData.GitHash = response.data.data.GitHead 
                        this.getData.GitEmail = response.data.data.CommitEmail;
                        this.getData.ServerPath = response.data.data.DeployPath
                        // 发布文件
                        this.getData.GitFileList = ""
                        for (var i=0;i<response.data.data.DeployFileList.length;i++){
                            this.getData.GitFileList += response.data.data.DeployFileList[i]+"\n"
                        }
                        // 目标服务器发布目录
                        this.IPList = response.data.data.DeployIP
                        this.getData.ServeripPath = ""
                        for(var i=0;i<this.IPList.length;i++){
                            this.getData.ServeripPath += this.IPList[i]+":"+this.getData.ServerPath+"    "
                        }
                    }else{
                        this.getData ={}
                        this.$set(this.getData, "GitInfo", "暂无内容");
                        this.$set(this.getData, "GitUrl",  "暂无内容");
                        this.$set(this.getData, "GitEmail", "暂无内容");
                        this.$set(this.getData, "GitHash", "暂无内容");
                        this.$set(this.getData, "ServeripPath", "暂无内容");
                        this.$message({
                            type: 'error',
                            message: response.data.message,
                            showClose: true, 
                            duration: 0,
                        });
                    }
                    return
                })   
            }
        },

        // 发布代码加载动画
        openFullScreen() {
            loading = this.$loading({
            lock: true,
            text: '项目发布中,请稍后...',
            spinner: 'el-icon-loading',
            background: 'rgba(0, 0, 0, 0.7)'
            })
        },

        // 发布代码
        deployCode(){
            this.sendData.deployment_id = this.values                       // 项目ID
            this.sendData.deployment_file_list = this.getData.GitFileList   // 发布文件
            this.sendData.deployment_commit = this.getData.GitInfo          // 备注
            this.sendData.git_head = this.getData.GitHash                   // 指针
            // 需要先拉取发布项目才能发布
            if (this.sendData.deployment_id == "") {
                this.$message({
                    type:'warning',
                    message:'请选择发布项目'
                })
                return
            }

            this.openFullScreen()

            this.$axios.post(config.urlpath+'/release/add', this.sendData).then(response => {           
                if (response.data.status == 200) {
                    // 发布成功
                    loading.close();
                    this.getLogList()
                    this.$message({message: '项目发布成功',showClose: true, duration:0, type: 'success'});
                    this.$set(this.getData, "GitInfo", "暂无内容");
                    this.$set(this.getData, "GitUrl",  "暂无内容");
                    this.$set(this.getData, "GitEmail", "暂无内容");
                    this.$set(this.getData, "GitHash", "暂无内容");
                    this.$set(this.getData, "ServeripPath", "暂无内容");
                    this.$set(this.getData, "GitFileList", "");
                    this.values = ""
                    return
                }
                // 发布失败
                loading.close();
                this.$message({message:response.data.message, showClose: true, duration:0, type: 'error'});
                this.getLogList()
                return
                
            }).catch(error => {
                this.$message({
                    duration:0,
                    showClose: true,
                    message: '项目发布失败',
                    type: 'error'
                });
                return
            })
        },
        
        // 回滚提示
        rollbackinfo(id){
            this.$confirm('确定回滚该项目吗?', '回滚项目', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(()=>{
                this.rollback(id)
            }).catch(()=>{
                this.$message({
                    type:'info',
                    message:'回滚已经取消'
                })
                return
            })
        },
        
        // 回滚代码
        rollback(id){
            this.openFullScreen()
            this.$axios.post(config.urlpath+'/release/rollback/'+id).then(response => {
                if (response.data.status != 200) {
                    this.$alert(response.data.message, '回滚失败', {
                        confirmButtonText: '确定',
                    })
                    loading.close();
                    return
                }
                loading.close();
                this.getLogList()
                this.$message({message: '项目回滚成功',showClose: true, duration:0, type: 'success'});
            })
        },
    },

    // 页面加载后自动执行
    mounted() {
        this.getreleasesList()
        this.getLogList()
    },
}
</script>