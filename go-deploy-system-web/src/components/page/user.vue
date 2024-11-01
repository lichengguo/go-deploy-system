<template>
    <div>
        <div class="crumbs">
            <el-breadcrumb separator="/">
                <el-breadcrumb-item>
                    <i class="el-icon-lx-calendar"></i> 用户列表
                </el-breadcrumb-item>
            </el-breadcrumb>
        </div>

        <div class="container">
            <!-- 搜索 -->
            <div class="handle-box">
                <el-input v-model="query.user_name" placeholder="请输入用户名" class="handle-input mr10" style="width:150px"></el-input>
                <el-button type="primary" icon="el-icon-search" @click="searchUserList">搜索</el-button>
                <el-button type="success" icon="el-icon-plus"   @click="addVisible = true">添加用户</el-button>
                <el-button type="warning" icon="el-icon-refresh-left" @click="load">刷新</el-button>
            </div>

            <!-- 用户列表 -->
            <el-table :data="UserList" border class="table" ref="multipleTable" header-cell-class-name="table-header" >
                <el-table-column prop="ID" label="UID" align="center"></el-table-column>
                <el-table-column prop="user_name" label="用户名" align="center"></el-table-column>
                <el-table-column prop="role" label="角色" align="center">
                    <template slot-scope="scope">
                        <span  v-if="scope.row.role==1">管理用户</span>
                        <span  v-if="scope.row.role==2">普通用户</span>
                    </template>
                </el-table-column>
                <el-table-column prop="Department.department_name" label="所属部门" align="center"></el-table-column>
                <el-table-column prop="status" label="状态" align="center">
                    <template slot-scope="scope">
                        <span  v-if="scope.row.status==1">可用</span>
                        <span  v-if="scope.row.status==2">冻结</span>
                    </template>
                </el-table-column>
                <el-table-column prop="CreatedAt" label="创建时间" width="170" align="center">
                    <template slot-scope="scope">
                        {{dateFormat(scope.row.CreatedAt)}}
                    </template>
                </el-table-column>
                <el-table-column prop="UpdatedAt" label="修改时间" width="170" align="center">
                    <template slot-scope="scope">
                        {{dateFormat(scope.row.UpdatedAt)}}
                    </template>
                </el-table-column>
                <el-table-column label="操作" width="200" align="center">
                    <template slot-scope="scope">
                        <el-button type="text" icon="el-icon-edit" @click="handleEdit(scope.row)" >修改</el-button>
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
        <el-dialog title="编辑" :visible.sync="editVisible" width="30%">
            <el-form ref="rulesform" :model="form" :rules="rulesform" label-width="80px">
                <el-form-item label="UID" prop="ID">
                    <el-input :disabled="true" v-model="form.ID"></el-input>
                </el-form-item>
                <el-form-item prop="user_name" label="用户名">
                    <el-input v-model="form.user_name"></el-input>
                </el-form-item>
                <el-form-item prop="role" label="角色">
                    <el-select v-model="form.role" placeholder="请选择角色" class="handle-input mr10" style="width:150px">
                        <el-option v-for="item in role" :key="item.id" :label="item.name" :value="item.id"> </el-option>
                    </el-select>
                </el-form-item>
                <el-form-item prop="status" label="状态">
                    <el-select v-model="form.status" placeholder="请选择状态" class="handle-input mr10" style="width:150px">
                        <el-option v-for="item in status" :key="item.id" :label="item.name" :value="item.id"> </el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="所属部门" prop="department_id" >
                    <el-select v-model="form.department_id" placeholder="请选择部门" class="handle-input mr10" style="width:150px">
                        <el-option v-for="item in departmentList" :key="item.ID" :label="item.department_name" :value="item.ID"></el-option>
                    </el-select>
                </el-form-item>
                <el-form-item prop="password" label="密码" v-if="false">
                    <el-input  type="password" v-model="form.password"></el-input>
                </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="editVisible = false">取 消</el-button>
                <el-button type="primary" @click="saveEdit">修 改</el-button>
            </span>
        </el-dialog>

        <!-- 添加弹出框 -->
        <el-dialog title="添加用户" :visible.sync="addVisible" width="30%">
            <el-form ref="rulesAddform" :model="addform" :rules="rulesAddform" label-width="80px">
                <el-form-item prop="user_name" label="用户名">
                    <el-input v-model="addform.user_name"></el-input>
                </el-form-item>
                <el-form-item prop="password" label="密码" >
                    <el-input  type="password" v-model="addform.password"></el-input>
                </el-form-item>
                <el-form-item prop="role" label="角色">
                    <el-select v-model="addform.role" placeholder="请选择角色" class="handle-input mr10" style="width:150px">
                        <el-option v-for="item in role" :key="item.id" :label="item.name" :value="item.id"></el-option>
                    </el-select>
                </el-form-item>
                <el-form-item prop="status" label="状态">
                    <el-select v-model="addform.status" placeholder="请选择状态" class="handle-input mr10" style="width:150px">
                        <el-option v-for="item in status" :key="item.id" :label="item.name" :value="item.id"></el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="所属部门" prop="department_id" >
                    <el-select v-model="addform.department_id" placeholder="请选择部门" class="handle-input mr10" style="width:150px">
                        <el-option v-for="item in departmentList" :key="item.ID" :label="item.department_name" :value="item.ID"></el-option>
                    </el-select>
                </el-form-item>
                
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="addVisible = false">取 消</el-button>
                <el-button type="primary" @click="addUser">添 加</el-button>
            </span>
        </el-dialog>
    </div>
</template>

<script>
import config from '../common/config.vue';

export default {
    data(){
        return {
            UserList:[],                    // 用户列表数据绑定
            departmentList:[],              // 部门数据列表
            editVisible: false,             // 编辑用户弹框
            form: {},                       // 编辑数据绑定
            addVisible: false,              // 添加用户弹框
            // 添加用户绑定数据
            addform:{
                user_name: '',
                password:'',
                role:2,
                department_id:'',
                status:1,
            },
            // 添加用户数据验证
            rulesAddform:{
                user_name:[{ required: true, message: '请输入英文名', trigger: 'blur'}],
                password:[{ required: true, message: '请输入密码', trigger: 'blur'}],
                role:[{ required: true, message: '请选择角色', trigger: 'blur'}],
                status:[{ required: true, message: '请选择状态', trigger: 'blur'}],
                department_id:[{ required: true, message: '请选择所属部门', trigger: 'blur'}],
            },
            // 编辑数据验证规则
            rulesform:{
                user_name:[{ required: true, message: '请输入英文名', trigger: 'blur'}],
                role:[{ required: true, message: '请选择角色', trigger: 'blur'}],
                status:[{ required: true, message: '请选择状态', trigger: 'blur'}],
                department_id:[{ required: true, message: '请选择所属部门', trigger: 'blur'}],
            },
            
            role:[{"id":1, "name":"管理用户"}, {"id":2,"name":"普通用户"}],  // 权限
            status:[{"id":1, "name":"使用中"}, {"id":2,"name":"已冻结"}],   // 状态

            // 日志分页
            pageTotal: 0,               // 用户列表总数
            
            // 搜索
            query: {
                user_name: '',          // 搜索部门名称
                pageIndex: 1,           // 搜索页码
                pageSize: 10,           // 搜索每页数量
            },
            
        };
    },

    methods: {
        // 添加用户
        addUser() {
            this.$refs.rulesAddform.validate((valid) => {
                if (!valid) {
                    return false
                }
                this.$axios.post(config.urlpath+'/user',this.addform).then(response =>{
                    if(response.data.status == 200){
                        this.$message.success('添加成功');
                        this.searchUserList()
                        this.addform = {
                            user_name: '',
                            password:'',
                            user_id:'',
                            role:2,
                            status:1,
                        };
                        this.addVisible=false;
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
                this.$axios.delete(config.urlpath+'/user/'+row.ID).then(response => {
                    if(response.data.status == 200){
                        this.$message.success('删除成功');
                        this.UserList.splice(index, 1);
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
            this.form.role = parseInt(this.form.role);
            this.editVisible = true;
        },

        // 保存编辑
        saveEdit() {
            this.$refs.rulesform.validate((valid) => {
                if (!valid) {
                    return false
                }
                this.$axios.put(config.urlpath+'/user/'+this.form.ID,this.form).then(response => {
                    if(response.data.status == 200){
                        this.$message.success(`修改成功`);
                        this.editVisible = false;
                        // this.load();
                        this.searchUserList()
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

        // 分页
        handlePageChange(val) {
            this.$set(this.query, 'pageIndex', val);
            this.searchUserList()   
        },

        // 刷新
        load() {
            this.query={
                user_name: '',
                pageIndex: 1,
                pageSize: 10,
            }
            this.searchUserList();
        },

        // 获取用户列表、搜索用户
        searchUserList(){
            this.$axios.get(config.urlpath+'/user', {
                params:{
                    pagesize: this.query.pageSize,
                    page: this.query.pageIndex,
                    user_name: this.query.user_name
                }
            }).then(response => {
                if(response.data.status == 200){
                    this.UserList = response.data.data;
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

        // 获取部门列表
        getDepartmentList() {
            this.$axios.get(config.urlpath+'/department').then(response => {
                if(response.data.status == 200){
                    this.departmentList = response.data.data;
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

        // 日期转换
        dateFormat(datetime){
            var date = new Date(datetime*1000);//时间戳为10位需*1000，时间戳为13位的话不需乘1000
            var year = date.getFullYear(),
                month = ("0" + (date.getMonth() + 1)).slice(-2),
                sdate = ("0" + date.getDate()).slice(-2),
                hour = ("0" + date.getHours()).slice(-2),
                minute = ("0" + date.getMinutes()).slice(-2),
                second = ("0" + date.getSeconds()).slice(-2);
            var result = year + "-"+ month +"-"+ sdate +" "+ hour +":"+ minute +":" + second;
            return result;
        },
    },

    // Vue生命周期函数
    created() {
        this.searchUserList();
        this.getDepartmentList();
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

.table-td-thumb {
    display: block;
    margin: auto;
    width: 40px;
    height: 40px;
}
</style>
