<template>
    <div>
        <div class="crumbs">
            <el-breadcrumb separator="/">
                <el-breadcrumb-item>
                    <i class="el-icon-lx-calendar"></i> 服务器列表
                </el-breadcrumb-item>
            </el-breadcrumb>
        </div>

        <div class="container">
            <!-- 搜索 -->
            <div class="handle-box">
                <el-input v-model="query.server_name" placeholder="服务器名称" class="handle-input mr10" style="width:150px"></el-input>
                <el-button type="primary" icon="el-icon-search" @click="searchServerList">搜索</el-button>
                <el-button type="success" icon="el-icon-plus"   @click="addVisible = true">添加服务器</el-button>
                <el-button type="warning" icon="el-icon-refresh-left" @click="load">刷新</el-button>
            </div>

            <!-- 服务器列表 -->
            <el-table :data="serverList" border ref="multipleTable" header-cell-class-name="table-header" >
                <el-table-column prop="ID" label="ID" align="center"></el-table-column>
                <el-table-column prop="Engineroom.engineroom_name" label="机房名称" ></el-table-column>
                <el-table-column prop="server_name" label="服务器名称"></el-table-column>
                <el-table-column prop="server_ip" label="服务器IP"></el-table-column>
                <el-table-column prop="server_port" label="服务器端口"></el-table-column>
                <el-table-column prop="server_user" label="服务器用户"></el-table-column>
                <el-table-column prop="server_pwd" label="服务器密码"></el-table-column>
                <el-table-column prop="server_key" label="服务器秘钥"></el-table-column>
                <el-table-column prop="server_status" label="服务器状态">
                    <template slot-scope="scope">
                        <span  v-if="scope.row.server_status==1">可使用</span>
                        <span  v-if="scope.row.server_status==2">冻结</span>
                    </template>
                </el-table-column>
                <el-table-column prop="CreatedAt" label="创建时间">
                    <template slot-scope="scope">
                        {{dateFormat(scope.row.CreatedAt)}}
                    </template>
                </el-table-column>
                <el-table-column prop="UpdatedAt" label="修改时间">
                    <template slot-scope="scope">
                        {{dateFormat(scope.row.UpdatedAt)}}
                    </template>
                </el-table-column>
                <el-table-column label="操作" width="180px" align="center">
                    <template slot-scope="scope">
                        <el-button type="text" icon="el-icon-refresh-left" style="color: #a6a9ad" @click="connectServer(scope.row)">测试</el-button>
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
                    @current-change="handlePageChange"
                ></el-pagination>
            </div>
        </div>

        <!-- 编辑弹出框 -->
        <el-dialog title="编辑" :visible.sync="editVisible" width="35%" top="0">
            <el-form ref="rulesform" :model="form" :rules="rulesform" label-width="110px">
                <el-form-item label="ID" prop="ID">
                    <el-input :disabled="true" v-model="form.ID"></el-input>
                </el-form-item>
                <el-form-item label="服务器名称" prop="server_name">
                    <el-input v-model="form.server_name"></el-input>
                </el-form-item>
                <el-form-item label="服务器名称IP" prop="server_ip">
                    <el-input v-model="form.server_ip"></el-input>
                </el-form-item>
                <el-form-item label="服务器端口" prop="server_port">
                    <el-input v-model="form.server_port"></el-input>
                </el-form-item>
                <el-form-item label="服务器用户" prop="server_user">
                    <el-input v-model="form.server_user"></el-input>
                </el-form-item>
                <el-form-item>
                    <el-radio v-model="radio" label="1" border>密码</el-radio>
                    <el-radio v-model="radio" label="2" border>秘钥</el-radio>
                </el-form-item>
                <template v-if="radio == '1'">
                    <el-form-item label="服务器密码" prop="server_pwd">
                        <el-input type="password" v-model="form.server_pwd"></el-input>
                    </el-form-item>
                </template>
                <template v-if="radio == '2'">
                    <el-form-item label="服务器秘钥" prop="server_key">
                        <el-input  v-model="form.server_key"></el-input>
                        <el-upload :headers="importHeaders" :action="UploadUrl()" :on-success="handleChange" name="keyfile" :file-list="fileList" :limit="1">
                            <i class="el-icon-upload"></i>
                            <div class="el-upload__text"><em>点击上传</em></div>
                        </el-upload>
                    </el-form-item>
                </template>
                <el-form-item label="所属机房" prop="engineroom_id" >
                    <el-select v-model="form.engineroom_id" placeholder="请选择机房" class="handle-input mr10" style="width:150px">
                        <el-option v-for="item in engineroomList" :key="item.ID" :label="item.engineroom_name" :value="item.ID"></el-option>
                    </el-select>
                </el-form-item>
                <el-form-item prop="server_status" label="状态">
                    <el-select v-model="form.server_status" placeholder="请选择状态" class="handle-input mr10" style="width:150px">
                        <el-option v-for="item in status" :key="item.id" :label="item.name" :value="item.id"></el-option>
                    </el-select>
                </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="editVisible = false">取 消</el-button>
                <el-button type="primary" @click="saveEdit">修 改</el-button>
            </span>
        </el-dialog>

        <!-- 添加弹出框 -->
        <el-dialog title="添加服务器" :visible.sync="addVisible" width="35%" top="0">
            <el-form ref="rulesform" :model="addform" :rules="rulesform" label-width="110px">
                <el-form-item label="服务器名称" prop="server_name">
                    <el-input v-model="addform.server_name"></el-input>
                </el-form-item>
                <el-form-item label="服务器名称IP" prop="server_ip">
                    <el-input v-model="addform.server_ip"></el-input>
                </el-form-item>
                <el-form-item label="服务器端口" prop="server_port">
                    <el-input v-model="addform.server_port"></el-input>
                </el-form-item>
                <el-form-item label="服务器用户名" prop="server_user">
                    <el-input v-model="addform.server_user"></el-input>
                </el-form-item>
                <el-form-item>
                    <el-radio v-model="radio" label="1" border>密码</el-radio>
                    <el-radio v-model="radio" label="2" border>秘钥</el-radio>
                </el-form-item>
                <template v-if="radio == '1'">
                    <el-form-item label="服务器密码" >
                        <el-input type="password" v-model="addform.server_pwd"></el-input>
                    </el-form-item>
                </template>               
                <template v-if="radio == '2'">
                    <el-form-item label="服务器秘钥" >
                        <el-upload :headers="importHeaders" :action="UploadUrl()" :on-success="handleChange" name="keyfile" :file-list="fileList" :limit="1">
                            <i class="el-icon-upload"></i>
                            <div class="el-upload__text"><em>点击上传</em></div>
                        </el-upload>
                    </el-form-item>
                </template>              
                <el-form-item label="所属机房" prop="engineroom_id" >
                    <el-select v-model="addform.engineroom_id" placeholder="请选择机房" class="handle-input mr10" style="width:150px">
                        <el-option v-for="item in engineroomList" :key="item.ID" :label="item.engineroom_name" :value="item.ID"></el-option>
                    </el-select>
                </el-form-item>
                <el-form-item prop="server_status" label="状态">
                    <el-select v-model="addform.server_status" placeholder="请选择状态" class="handle-input mr10" style="width:150px">
                        <el-option v-for="item in status" :key="item.id" :label="item.name" :value="item.id"></el-option>
                    </el-select>
                </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="addVisible = false">取 消</el-button>
                <el-button type="primary" @click="addServer">添 加</el-button>
            </span>
        </el-dialog>
    </div>
</template>

<script>
import config from '../common/config.vue';

export default {
    data() {
        return {
            addVisible: false,  // 添加服务器弹窗
            // 添加服务器绑定数据
            addform:{
                server_name: '',
                server_ip:'',
                server_port:'',
                engineroom_id:'',
                server_user:'',
                server_pwd:'',
                server_key:'',
                server_status:1,
            },
            // 添加服务器、修改服务器 数据校验
            rulesform:{
                server_name:[{ required: true, message: '请输入服务器名称', trigger: 'blur'}],
                server_ip:[{ required: true, message: '请输入服务器IP', trigger: 'blur'}],
                server_port:[{ required: true, message: '请输入服务器端口', trigger: 'blur'}],
                engineroom_id:[{ required: true, message: '请选择所属机房', trigger: 'blur'}],
                server_user:[{ required: true, message: '请输入服务器用户', trigger: 'blur'}],
                server_status:[{ required: true, message: '请选择状态', trigger: 'blur'}],
            },
            radio:"1",              // 添加服务器和修改服务器 密码、秘钥选择标志位
            fileList:[],            // 服务器秘钥
            // 上传服务器秘钥时；携带token
            importHeaders:{
                Authorization: `Bearer ${localStorage.getItem('token')}`, // token
            },

            status:[{"id":1, "name":"可使用"}, {"id":2,"name":"冻结"}],    // 服务器状态
            
            editVisible: false,     // 修改服务器弹窗
            form: {},               // 修改服务器数据绑定

            serverList:[],         // 服务器列表
            engineroomList:[],      // 机房列表
            
            // 搜索
            query: {
                server_name: '',
                pageIndex: 1,
                pageSize: 10,
            },
            
            // 分页
            pageTotal: 0,           
        };
    },

    methods: {
        // 添加操作
        addServer() {
            this.$refs.rulesform.validate((valid) => {
                if (!valid) {
                    return false
                }
                this.$axios.post(config.urlpath+'/server',this.addform).then(response => {
                    if(response.data.status == 200){
                        this.$message.success('添加成功');
                        this.searchServerList();
                        this.addform = {
                            server_name: '',
                            server_ip:'',
                            server_port:'',
                            engineroom_id:'',
                            server_user:'',
                            server_pwd:'',
                            server_key:'',
                            server_status:1,
                        };
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

        // 秘钥上传链接
        UploadUrl:function(){
            return config.urlpath+'/upload';   
        },
        // 秘钥上传成功以后
        handleChange(res, file){
            if(file.response.status==200){
                this.addform.server_key = res.url
                if(this.form){
                    this.form.server_key = res.url
                }           
            }else{
                this.$message.error(file.response.msg);
            }
        },

        // 获取机房列表
        getEngineroom(){
            this.$axios.get(config.urlpath+'/engineroom').then(response =>{
                if(response.data.status == 200){
                    this.engineroomList = response.data.data;
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

        // 删除操作
        handleDelete(index,row) {
            // 二次确认删除
            this.$confirm('确定要删除吗？', '提示', {
                type: 'warning'
            }).then(() => {
                this.$axios.delete(config.urlpath+'/server/'+row.ID).then(response =>{
                    if(response.data.status == 200){
                        this.$message.success('删除成功');
                        this.serverList.splice(index, 1);
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
            this.form.ID = parseInt(this.form.ID);
            this.editVisible = true;
        },

        // 保存编辑
        saveEdit() {
            this.$refs.rulesform.validate((valid) => {
                if (!valid) {
                    return false
                }
                this.$axios.put(config.urlpath+'/server/'+this.form.ID,this.form).then(response =>{
                    if(response.data.status == 200){
                        this.$message.success(`修改成功`);
                        this.editVisible = false;
                        this.searchServerList()
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

        // 日期转换
        dateFormat(datetime){
            var date = new Date(datetime*1000);//时间戳为10位需*1000，时间戳为13位的话不需乘1000
            var year = date.getFullYear(),
                month = ("0" + (date.getMonth() + 1)).slice(-2),
                sdate = ("0" + date.getDate()).slice(-2),
                hour = ("0" + date.getHours()).slice(-2),
                minute = ("0" + date.getMinutes()).slice(-2),
                second = ("0" + date.getSeconds()).slice(-2);
            // 拼接
            var result = year + "-"+ month +"-"+ sdate +" "+ hour +":"+ minute +":" + second;
            // 返回
            return result;
        },

        // 分页
        handlePageChange(val) {
            this.$set(this.query, 'pageIndex', val);
            this.searchServerList()
        },

        // 刷新
        load(){
            this.query={
                server_name: '',
                pageIndex: 1,
                pageSize: 10,
            }
            this.searchServerList()
        },

        // 获取服务器列表、搜索服务器
        searchServerList(){
            this.$axios.get(config.urlpath+'/server', {
                params:{
                    pagesize: this.query.pageSize,
                    page: this.query.pageIndex,
                    server_name: this.query.server_name
                }
            }).then(response => {
                if(response.data.status == 200){
                    this.serverList = response.data.data;
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

        // 服务器测试连接
        connectServer(row){
            this.$axios.get(config.urlpath + '/server/connect/' + row.ID).then(response => {
                if (response.data.status == 200){
                    this.$message({
                        type: 'success',
                        message: '服务器测试连接成功'
                    })
                } else {
                    this.$message({
                        type: 'error',
                        message: '服务器测试连接失败'
                    })
                }
            }).catch(error => {
                this.$message({
                    type: 'error',
                    message: error
                })
                return
            })
        },
        
    },

    // vue生命周期
    created() {
        this.searchServerList()
        this.getEngineroom()
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
