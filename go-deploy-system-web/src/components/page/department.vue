<template>
    <div>
        <div class="crumbs">
            <el-breadcrumb separator="/">
                <el-breadcrumb-item>
                    <i class="el-icon-lx-calendar"></i> 部门列表
                </el-breadcrumb-item>
            </el-breadcrumb>
        </div>
        
        <div class="container">
            <!-- 搜索 -->
            <div class="handle-box">
                <el-input v-model="query.department_name" placeholder="部门名称" class="handle-input mr10" style="width:150px"></el-input>
                <el-button type="primary" icon="el-icon-search" @click="searchDepartmentList">搜索</el-button>
                <el-button type="success" icon="el-icon-plus"   @click="addVisible = true">添加部门</el-button>
                <el-button type="warning" icon="el-icon-refresh-left" @click="load">刷新</el-button>
            </div>
            
            <!-- 部门列表 -->
            <el-table :data="departmentList" border class="table" ref="multipleTable" header-cell-class-name="table-header">
                <el-table-column prop="ID" label="ID"  align="center"></el-table-column>
                <el-table-column prop="department_name" label="部门名称" align="center"></el-table-column>
                <el-table-column prop="CreatedAt" label="创建时间" align="center">
                    <template slot-scope="scope">
                        {{dateFormat(scope.row.CreatedAt)}}
                    </template>
                </el-table-column>
                <el-table-column prop="UpdatedAt" label="修改时间" align="center" >
                    <template slot-scope="scope">
                        {{dateFormat(scope.row.UpdatedAt)}}
                    </template>
                </el-table-column>
                <el-table-column label="操作" align="center">
                    <template slot-scope="scope">
                        <el-button type="text" icon="el-icon-edit" @click="handleEdit(scope.row)" >修改</el-button>
                        <el-button type="text" icon="el-icon-delete" class="red" @click="handleDelete(scope.$index,scope.row)">删除</el-button>
                    </template>
                </el-table-column>
            </el-table>
            
            <!-- 分页 -->
            <div class="pagination">
                <el-pagination 
                    small
                    background 
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
                <el-form-item label="ID" prop="ID">
                    <el-input :disabled="true" v-model="form.ID"></el-input>
                </el-form-item>
                <el-form-item label="部门名称" prop="department_name">
                    <el-input v-model="form.department_name"></el-input>
                </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="editVisible = false">取 消</el-button>
                <el-button type="primary" @click="saveEdit">修 改</el-button>
            </span>
        </el-dialog>

        <!-- 添加弹出框 -->
        <el-dialog title="添加部门" :visible.sync="addVisible" width="30%">
            <el-form ref="rulesAddform" :model="addform" :rules="rulesAddform" label-width="80px">
                <el-form-item label="部门名称" prop="department_name">
                    <el-input v-model="addform.department_name"></el-input>
                </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="addVisible = false">取 消</el-button>
                <el-button type="primary" @click="addDpartment">添 加</el-button>
            </span>
        </el-dialog>
    </div>
</template>

<script>
import config from '../common/config.vue';

export default {
    data() {
        return {
            departmentList:[],          // 部门列表显示数据
            pageTotal: 0,               // 部门列表总数
            editVisible: false,         // 编辑部门弹框
            form: {},                   // 编辑部门表单数据
            addVisible: false,          // 添加部门弹框    
            addform:{                   // 添加部门表单数据
                department_name: '',
                user_id:''
            },

            // 搜索
            query: {
                department_name: '',    // 搜索部门名称
                pageIndex: 1,           // 搜索页码
                pageSize: 10,           // 搜索每页数量
            },
            
            // 校验规则
            // 编辑
            rulesform:{                 
                department_name:[{ required: true, message: '请输入部门名称', trigger: 'blur'}],
            },
            // 添加
            rulesAddform:{
                department_name:[{ required: true, message: '请输入部门名称', trigger: 'blur'}],
            },
        };
    },

    methods: {
        // 添加部门
        addDpartment() {
            this.$refs.rulesAddform.validate((valid) => {
                if (!valid) {
                    return false
                }
                this.$axios.post(config.urlpath+'/department',this.addform).then(response =>{
                    if(response.data.status == 200){
                        this.$message.success('添加成功');
                        this.searchDepartmentList()
                        this.addform = {
                            department_name: '',
                            user_id:''
                        };
                        this.addVisible=false;
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
            })
        },

        // 删除操作
        handleDelete(index,row) {
            // 二次确认删除
            this.$confirm('确定要删除吗？', '提示', {
                type: 'warning'
            }).then(() => {
                this.$axios.delete(config.urlpath+'/department/'+row.ID).then(response => {
                    if(response.data.status == 200){
                        this.$message.success('删除成功');
                        this.departmentList.splice(index, 1);
                    }else{
                        this.$message({
                            type: 'error',
                            message: response.data.massage
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
                this.$axios.put(config.urlpath+'/department/'+this.form.ID,this.form).then(response =>{
                    if(response.data.status == 200){
                        this.$message.success(`修改成功`);
                        this.editVisible = false;
                        this.searchDepartmentList()
                    }else{
                        this.$message({
                            type: 'error',
                            message: response.data.message
                        })
                    }
                    return
                }).catch(error =>{
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
            this.searchDepartmentList()
        },
        
        // 刷新
        load(){
            this.query.department_name = ''
            this.query.pageIndex = 1
            this.query.pageSize = 10
            this.searchDepartmentList()
        },

        // 获取部门列表、搜索
        searchDepartmentList(){
            this.$axios.get(config.urlpath+'/department', {
                params:{
                    pagesize: this.query.pageSize,
                    page: this.query.pageIndex,
                    department_name: this.query.department_name
                }
            }).then(response => {
                if(response.data.status == 200){
                    this.departmentList = response.data.data;
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
    },

    // Vue生命周期函数
    created() {
        this.searchDepartmentList()
    },
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
