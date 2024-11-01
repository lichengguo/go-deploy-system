<template>
    <div>
        <div class="crumbs">
            <el-breadcrumb separator="/">
                <el-breadcrumb-item>
                    <i class="el-icon-lx-calendar"></i> 发布日志
                </el-breadcrumb-item>
            </el-breadcrumb>
        </div>
        
        <el-card>
            <div style="margin-bottom:15px">
                <el-input v-model="query.deployment_name"  placeholder="请输入项目名称" style="width:150px"></el-input>
                <el-input v-show="query.showInput" v-model="query.deployment_user_name"  placeholder="请输入发布人名称" style="width:150px"></el-input>
                <el-button type="primary" icon="el-icon-search" @click="searchDeploymentLog">搜索</el-button>
                <el-button type="warning" icon="el-icon-refresh-left" @click="load">刷新</el-button>
            </div>
            <el-row>
                <!-- 日志展示 -->
                <el-col :span="24">
                    <el-table :data="logList" border  header-cell-class-name="table-header" >
                        <el-table-column prop="ID" label="ID" width="80"></el-table-column>
                        <el-table-column prop="deployment_id" width="80" label="项目ID" align="center"></el-table-column>
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
                    
                    <!-- 分页 -->
                    <div class="pagination">
                        <el-pagination
                            background
                            small
                            layout="total, prev, pager, next"
                            :current-page="query.pageIndex"
                            :page-size="query.pageSize"
                            :total="pageTotal"
                            @current-change="handlePageChange"
                        ></el-pagination>
                    </div>
                </el-col>
            </el-row>
        </el-card>
        
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

            pageTotal:0,                    // 数据总量
            
            // 搜索
            query: {
                showInput:false,            // 发布人名称搜索框是否隐藏
                deployment_name:'',         // 项目名称
                deployment_user_name:'',    // 发布人
                pageIndex: 1,               // 页码
                pageSize: 10,               // 分页数据量
            },
        }
    },

    methods:{
        // 发布代码加载动画
        openFullScreen() {
            loading = this.$loading({
            lock: true,
            text: '项目发布中,请稍后...',
            spinner: 'el-icon-loading',
            background: 'rgba(0, 0, 0, 0.7)'
            })
        },

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
            this.form.deployment_status = row.deployment_status
            if (row.deployment_status == 1) {
                this.form.deployment_status = '成功'
            } else {
                this.form.deployment_status = '失败'
            }
            this.form.deployment_fail_info = row.deployment_fail_info
            this.form.git_head = row.git_head
            this.form.UpdatedAt = this.dateFormat(row.UpdatedAt)
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
                this.searchDeploymentLog()
                this.$message({message: '项目回滚成功',showClose: true, duration:0, type: 'success'});
            })
        },

        // 分页
        handlePageChange(val) {
            this.$set(this.query, 'pageIndex', val);
            this.searchDeploymentLog() 
        },

        // 刷新日志
        load(){
            this.query.deployment_user_name = ''
            this.query.deployment_name = ''
            this.query.pageIndex = 1
            this.query.pageSize = 10
            this.searchDeploymentLog()
        },

        // 获取日志列表 日志搜索
        searchDeploymentLog(){
            this.$axios.get(config.urlpath+'/deploymentlogs', {
                params:{
                    page: this.query.pageIndex,
                    pagesize: this.query.pageSize,
                    deployment_name: this.query.deployment_name,
                    deployment_user_name: this.query.deployment_user_name,
                },
            }).then(response =>{
                if (response.data.status == 200) {
                    this.pageTotal = response.data.total
                    this.logList = response.data.data
                } else {
                     alert('获取日志失败')
                }
            })
        },
    },

    // 页面加载后自动执行
    mounted() {
        this.searchDeploymentLog()

        // 管理员才可以使用 发布人 条件搜索发布日志
        var role = localStorage.getItem('role')
        if (role == 1) {
            this.query.showInput = true
        }
        
    },
}
</script>