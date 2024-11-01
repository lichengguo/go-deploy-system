<template>
    <div>
        <div class="crumbs">
            <el-breadcrumb separator="/">
                <el-breadcrumb-item>
                    <i class="el-icon-lx-calendar"></i> 项目配置列表
                </el-breadcrumb-item>
            </el-breadcrumb>
        </div>

        <div class="container">
            <!-- 搜索 -->
            <div class="handle-box">
                <el-input v-model="query.deployment_name" placeholder="项目名称" class="handle-input mr10" style="width:150px"></el-input>
                <el-button type="primary" icon="el-icon-search" @click="searchProjectList">搜索</el-button>
                <el-button type="success" icon="el-icon-plus"   @click="addVisible = true">添加项目配置</el-button>
                <el-button type="warning" icon="el-icon-refresh-left" @click="load">刷新</el-button>
            </div>

            <!-- 项目配置列表 -->
            <el-table :data="projectList" border ref="multipleTable" header-cell-class-name="table-header">
                <el-table-column prop="ID" label="ID" align="center"></el-table-column>
                <el-table-column prop="deploy_name" label="项目名称"></el-table-column>
                <el-table-column prop="git_url_http" label="Git URL"></el-table-column>
                <el-table-column prop="git_url_ssh" label="Git SSH"></el-table-column>
                <el-table-column prop="git_branch" label="Git分支"></el-table-column>
                <el-table-column prop="git_user" label="Git用户名"></el-table-column>
                <el-table-column prop="git_passwd" label="Git密码"></el-table-column>
                <el-table-column prop="git_key" label="Git密钥地址"></el-table-column>
                <el-table-column prop="deploy_server_path" label="服务器发布目录"></el-table-column>
                <el-table-column prop="server_key" label="所属服务器"  width="180">
                     <template slot-scope="scope">
                        <div v-for="val in scope.row.ServerList" :key="val">
                            {{val.Engineroom.engineroom_name + "-" + val.server_name}}
                        </div>
                    </template>
                </el-table-column>
                <el-table-column prop="user_key" label="发布人">
                    <template slot-scope="scope">
                        <div v-for="val in scope.row.UserList" :key="val">{{val.user_name}}</div>
                    </template>
                </el-table-column>
                <el-table-column label="操作" width="170" align="center">
                    <template slot-scope="scope">
                        <el-button type="text" icon="el-icon-edit" @click="handleEdit(scope.row)">修改</el-button>
                        <el-button type="text" icon="el-icon-delete" class="red" @click="handleDelete(scope.$index,scope.row)">删除</el-button>
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
                    @current-change="handlePageChange">
                </el-pagination>
            </div>
        </div>

        <!-- 编辑弹出框 -->
        <el-dialog title="编辑" :visible.sync="editVisible" width="38%" top="0">
            <el-form ref="rulesSaveform" :model="form" :rules="rulesSaveform" label-width="110px">
                <el-form-item label="ID" prop="ID">
                    <el-input :disabled="true" v-model="form.ID"></el-input>
                </el-form-item>
                <el-form-item label="项目名称" prop="deploy_name">
                    <el-input v-model="form.deploy_name"></el-input>
                </el-form-item>
                <el-form-item label="Git分支" prop="git_branch">
                    <el-input v-model="form.git_branch"></el-input>
                </el-form-item>
                <el-form-item>
                    <el-radio v-model="radio" label="1" border>SSH</el-radio>
                    <el-radio v-model="radio" label="2" border>URL</el-radio>
                </el-form-item>
                <template v-if="radio == '2'">
                    <el-form-item label="Git URL" prop="git_url_http">
                    <el-input v-model="form.git_url_http"></el-input>
                    </el-form-item>
                    <el-form-item label="Git用户名" prop="git_user">
                        <el-input v-model="form.git_user"></el-input>
                    </el-form-item>
                    <el-form-item label="Git密码" prop="git_passwd">
                        <el-input type="password" v-model="form.git_passwd"></el-input>
                    </el-form-item>
                </template>
                <template v-if="radio == '1'">
                    <el-form-item label="Git SSH" prop="git_url_ssh">
                        <el-input v-model="form.git_url_ssh"></el-input>
                    </el-form-item>
                    <el-form-item label="Git密钥地址" prop="git_key">
                        <el-input v-model="form.git_key"></el-input>
                        <el-upload :headers="importHeaders" :action="UploadUrl()" :on-success="handleChange" name="keyfile" :file-list="fileList" :limit="1">
                            <i class="el-icon-upload"></i>
                            <div class="el-upload__text"><em>点击上传</em></div>
                        </el-upload>
                    </el-form-item>
                </template>
                <el-form-item label="发布目录" prop="deploy_server_path">
                    <el-input v-model="form.deploy_server_path"></el-input>
                </el-form-item>
                <el-form-item label="所属服务器">
                    <el-checkbox-group v-model="saveform.server_id" @change="changeFunc()">
                        <el-checkbox v-for="servers in serversList" :label="servers.ID" :key="servers.ID" >{{servers.Engineroom.engineroom_name + "-" +servers.server_name}}</el-checkbox>
                    </el-checkbox-group>
                </el-form-item>
                <el-form-item label="选择发布人">
                    <el-checkbox-group v-model="saveform.user_id" @change="changeFunc()">
                        <el-checkbox v-for="user in userList" :label="user.ID" :key="user.ID">{{user.user_name}}</el-checkbox>
                    </el-checkbox-group>
                </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="editVisible = false">取 消</el-button>
                <el-button type="primary" @click="saveEdit">修 改</el-button>
            </span>
        </el-dialog>

        <!-- 添加弹出框 -->
        <el-dialog title="添加项目" :visible.sync="addVisible" width="38%" top="0">
            <el-form ref="rulesAddform" :model="addform" :rules="rulesAddform" label-width="110px">
                <el-form-item label="项目名称" prop="deploy_name">
                    <el-input v-model="addform.deploy_name"></el-input>
                </el-form-item>
                <el-form-item label="Git分支" prop="git_branch">
                    <el-input v-model="addform.git_branch"></el-input>
                </el-form-item>
                <el-form-item>
                    <el-radio v-model="radio" label="1" border>SSH</el-radio>
                    <el-radio v-model="radio" label="2" border>URL</el-radio>
                </el-form-item>
                <template v-if="radio == '2'">
                    <el-form-item label="Git URL" >
                        <el-input v-model="addform.git_url_http"></el-input>
                    </el-form-item>
                    <el-form-item label="Git用户名" >
                        <el-input v-model="addform.git_user"></el-input>
                    </el-form-item>
                    <el-form-item label="Git密码" >
                        <el-input type="password" v-model="addform.git_passwd"></el-input>
                    </el-form-item>
                </template>
                <template v-if="radio == '1'">
                    <el-form-item label="Git SSH" >
                        <el-input v-model="addform.git_url_ssh"></el-input>
                    </el-form-item>
                    <el-form-item label="Git密钥地址" >
                    <el-upload :headers="importHeaders" :action="UploadUrl()" :on-success="handleChange" name="keyfile" :file-list="fileList" :limit="1">
                            <i class="el-icon-upload"></i>
                            <div class="el-upload__text"><em>点击上传</em></div>
                        </el-upload>
                    </el-form-item>
                </template>
                <el-form-item label="发布目录" prop="deploy_server_path">
                    <el-input v-model="addform.deploy_server_path"></el-input>
                </el-form-item>
                <el-form-item label="所属服务器" prop="server_id">
                    <el-checkbox-group v-model="addform.server_id">
                        <el-checkbox v-for="servers in serversList" :label="servers.ID" :key="servers.ID">{{servers.Engineroom.engineroom_name + "-" +servers.server_name}}</el-checkbox>
                    </el-checkbox-group>
                </el-form-item>
                <el-form-item label="选择发布人" prop="user_id">
                    <el-checkbox-group v-model="addform.user_id">
                        <el-checkbox v-for="user in userList" :label="user.ID" :key="user.ID">{{user.user_name}}</el-checkbox>
                    </el-checkbox-group>
                </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="addVisible = false">取 消</el-button>
                <el-button type="primary" @click="addProject">添 加</el-button>
            </span>
        </el-dialog>

    </div>
</template>

<script>
import config from '../common/config.vue';

export default {
    data() {
        return {
            projectList:[],     // 项目列表数据
            serversList:[],     // 服务器列表数据
            userList:[],        // 用户列表数据
            radio:"1",          // 标志位 添加项目、编辑项目时
            fileList:[],        // 上传文件列表 添加项目

            addVisible: false,  // 添加项目弹框
            // 添加项目绑定数据
            addform:{
                deploy_name: '',
                git_url_http:'',
                git_url_ssh:'',
                git_branch:'',
                git_user:'',
                git_passwd:'',
                git_key:'',
                server_id:[],
                deploy_server_path:'',
                user_id:[]
            },
            // 添加项目数据验证
            rulesAddform:{
                deploy_name:[{ required: true, message: '请输入发布项目名称', trigger: 'blur'}],
                git_branch:[{ required: true, message: '请输入Git分支', trigger: 'blur'}],
                user_id:[{ required: true, message: '请选择发布用户', trigger: 'blur'}],
                deploy_server_path:[{ required: true, message: '请输入发布目录', trigger: 'blur'}],
                server_id:[{ required: true, message: '请选择所属机房', trigger: 'blur'}],
                user_id:[{ required: true, message: '请选择发布人', trigger: 'blur'}],
            },

            // 上传秘钥时；携带token
            importHeaders:{
                Authorization: `Bearer ${localStorage.getItem('token')}`, // token
            },

            editVisible: false,     // 修改项目弹窗
            form: {},               // 修改项目数据绑定
            // 修改项目数据验证规则
            rulesSaveform:{
                deploy_name:[{ required: true, message: '请输入发布项目名称', trigger: 'blur'}],
                git_branch:[{ required: true, message: '请输入Git分支', trigger: 'blur'}],
                user_id:[{ required: true, message: '请选择发布用户', trigger: 'blur'}],
                deploy_server_path:[{ required: true, message: '请输入发布目录', trigger: 'blur'}],
            },
            // 修改项目时；存储服务器ID和用户ID
            saveform: {
                server_id:[],
                user_id:[],
            },
            
            // 搜索
            query: {
                deployment_name: '',
                pageIndex: 1,
                pageSize: 10,
            },
            
            // 分页
            pageTotal: 0,   // 数据总数
            
        };
    },
    methods: {
        // 获取服务器列表
        getServersList(){
            this.$axios.get(config.urlpath+'/server').then(response =>{
                if(response.data.status == 200){
                    this.serversList = response.data.data;
                }
                return
            }).catch(error => {
                this.$message({
                    type: 'error',
                    message: error
                })
                return
            });
        },

        // 获取用户列表
        getUserList(){
            this.$axios.get(config.urlpath+'/user').then(response => {
                if(response.data.status == 200){
                    this.userList = response.data.data;
                }
                return
            }).catch(error => {
                this.$message({
                    type: 'error',
                    message: error
                })
                return
            });
        },

        // 添加项目
        addProject() {   
            this.$refs.rulesAddform.validate((valid) => {
                if (!valid) {
                    return false
                }
                this.$axios.post(config.urlpath+'/deployment',this.addform).then(response=>{
                    if(response.data.status == 200){
                        this.$message.success('添加成功');
                        this.searchProjectList();
                        this.addform = {
                            deploy_name: '',
                            git_url_http:'',
                            git_url_ssh:'',
                            git_branch:'',
                            git_user:'',
                            git_passwd:'',
                            git_key:'',
                            server_id:[],
                            deploy_server_path:'',
                            user_id:[]
                        };
                        this.radio="1";
                        this.addVisible=false;
                        this.fileList = []
                    }else{
                        this.$message({
                            type: 'error',
                            message: response.data.message
                        })
                    }
                    return
                }).catch(error => {
                    this.$message({
                        type: 'error',
                        message: error
                    })
                    return
                });
            })
        },

        // 删除操作
        handleDelete(index,row) {
            // 二次确认删除
            this.$confirm('确定要删除吗？', '提示', {
                type: 'warning'
            }).then(() => {
                this.$axios.delete(config.urlpath+'/deployment/'+row.ID).then(response =>{
                    if(response.data.status == 200){
                        this.$message.success('删除成功');
                        this.projectList.splice(index, 1);
                    }else{
                        this.$message({
                            type: 'error',
                            message: response.data.massage
                        })
                    }
                    return
                }).catch(error => {
                    this.$message({
                        type: 'error',
                        message: error
                    })
                    return
                });
            }).catch(() => {});
        },

        // 编辑操作
        handleEdit(row) {
            this.form = row;
            
            var user_id = [];
            if(this.form.UserList){
                for(var i=0; i<this.form.UserList.length;i++){
                    user_id.push(this.form.UserList[i].ID)
                }
            }            
            
            var server_id = [];
            if(this.form.ServerList){
                for(var i=0; i<this.form.ServerList.length; i++){
                    server_id.push(this.form.ServerList[i].ID)
                }
            }    
            
            this.saveform.user_id = user_id;
            this.saveform.server_id = server_id;
            this.form.user_id = user_id;
            this.form.server_id = server_id;
            this.form.ID = parseInt(this.form.ID);

            this.editVisible = true;
            if(this.form.git_url_ssh){
                this.radio = "1";
            }else{
                this.radio = "2";
            }
        },

        // 保存编辑
        saveEdit() {
            this.$refs.rulesSaveform.validate((valid) => {
                if (!valid) {
                    return false
                }
                this.$axios.put(config.urlpath+'/deployment/'+this.form.ID,this.form).then(response =>{
                    if(response.data.status == 200){
                        this.$message.success(`修改成功`);
                        this.editVisible = false;
                        this.searchProjectList()
                        this.fileList = []
                    }else{
                        this.$message({
                            type: 'error',
                            message: response.data.message
                        })
                    }
                    return
                }).catch(error => {
                    this.$message({
                        type: 'error',
                        message: error
                    })
                    return
                });
            })
        },

        // 秘钥上传成功以后钩子函数
        handleChange(res, file){
            if(file.response.status==200){
                this.addform.git_key = res.url
                if(this.form){
                    this.form.git_key = res.url
                } 
            }else{
                this.$message.error(file.response.msg);
            }
        },
        // 上传链接
        UploadUrl:function(){
            return config.urlpath+'/upload';   
        },

        // 分页
        handlePageChange(val) {
            this.$set(this.query, 'pageIndex', val);
            this.searchProjectList()  
        },

        // 刷新
        load(){
            this.query={
                deployment_name: '',
                pageIndex: 1,
                pageSize: 10,
            }
            this.searchProjectList()
        },

        // 搜索项目
        searchProjectList(){
            this.$axios.get(config.urlpath+'/deployment', {
                params:{
                    pagesize: this.query.pageSize,
                    page: this.query.pageIndex,
                    deployment_name: this.query.deployment_name
                }
            }).then(response => {
                if(response.data.status == 200){
                    this.projectList = response.data.data;
                    this.pageTotal =  response.data.total;
                }else{
                    this.$message({
                        type: 'error',
                        message: response.data.message
                    })
                }
                return
            }).catch(error => {
                this.$message({
                    type: 'error',
                    message: error
                })
                return
            });
        },

        // 添加编辑时；服务器、用户选择函数
        changeFunc(){
            this.form.user_id = this.saveform.user_id
            this.form.server_id = this.saveform.server_id
        },

    },

    // vue生命周期
    created() {
        this.searchProjectList()
        this.getServersList()
        this.getUserList()
   }
};
</script>

<style scoped>
.handle-box {
    margin-bottom: 20px;
}

.handle-input {
    width: 300px;
    display: inline-block;
}

.red {
    color: #ff0000;
}

.mr10 {
    margin-right: 10px;
}
</style>
