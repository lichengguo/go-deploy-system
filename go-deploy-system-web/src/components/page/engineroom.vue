<template>
    <div>
        <div class="crumbs">
            <el-breadcrumb separator="/">
                <el-breadcrumb-item>
                    <i class="el-icon-lx-calendar"></i> 机房列表
                </el-breadcrumb-item>
            </el-breadcrumb>
        </div>
        <div class="container">
            <!-- 搜索 -->
            <div class="handle-box">
                <el-input v-model="query.engineroom_name" placeholder="机房名称" class="handle-input mr10" style="width:150px"></el-input>
                <el-button type="primary" icon="el-icon-search" @click="searchEngineroomList">搜索</el-button>
                <el-button type="success" icon="el-icon-plus"   @click="addVisible = true">添加机房</el-button>
                <el-button type="warning" icon="el-icon-refresh-left" @click="load">刷新</el-button>
            </div>

            <!-- 机房列表 -->
            <el-table :data="computerList" border class="table" ref="multipleTable" header-cell-class-name="table-header">
                <el-table-column prop="ID" label="ID" align="center"></el-table-column>
                <el-table-column prop="engineroom_name" label="机房名称" ></el-table-column>
                <el-table-column prop="contact" label="联系人"></el-table-column>
                <el-table-column prop="contact_info" label="联系方式"></el-table-column>
                <el-table-column prop="address" label="机房地址"></el-table-column>
                <el-table-column prop="CreatedAt" label="创建时间" width="170">
                    <template slot-scope="scope">
                        {{dateFormat(scope.row.UpdatedAt)}}
                    </template>
                </el-table-column>
                <el-table-column prop="UpdatedAt" label="修改时间" width="170">
                    <template slot-scope="scope">
                        {{dateFormat(scope.row.UpdatedAt)}}
                    </template>
                </el-table-column>
                <el-table-column label="操作" width="210" align="center">
                    <template slot-scope="scope">
                        <el-button type="text" icon="el-icon-edit" @click="handleEdit(scope.row)">修改</el-button>
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
                    @current-change="handlePageChange"
                ></el-pagination>
            </div>
        </div>

        <!-- 编辑弹出框 -->
        <el-dialog title="编辑" :visible.sync="editVisible" width="30%">
            <el-form ref="rulesform" :model="form" :rules="rulesform" label-width="80px">
                <el-form-item label="ID" prop="ID">
                    <el-input :disabled="true" v-model="form.ID"></el-input>
                </el-form-item>
                <el-form-item label="机房名称" prop="engineroom_name">
                    <el-input v-model="form.engineroom_name"></el-input>
                </el-form-item>
                <el-form-item label="联系人" prop="contact">
                    <el-input v-model="form.contact"></el-input>
                </el-form-item>
                <el-form-item label="联系方式" prop="contact_info">
                    <el-input v-model="form.contact_info"></el-input>
                </el-form-item>
                <el-form-item label="机房地址" prop="address">
                    <el-input v-model="form.address"></el-input>
                </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="editVisible = false">取 消</el-button>
                <el-button type="primary" @click="saveEdit">修 改</el-button>
            </span>
        </el-dialog>

        <!-- 添加弹出框 -->
        <el-dialog title="添加机房" :visible.sync="addVisible" width="30%">
            <el-form ref="rulesform" :model="addform" :rules="rulesform" label-width="80px">
                <el-form-item label="机房名称" prop="engineroom_name">
                    <el-input v-model="addform.engineroom_name"></el-input>
                </el-form-item>
                <el-form-item label="联系人" prop="contact">
                    <el-input v-model="addform.contact"></el-input>
                </el-form-item>
                <el-form-item label="联系方式" prop="contact_info">
                    <el-input v-model="addform.contact_info"></el-input>
                </el-form-item>
                <el-form-item label="机房地址" prop="address">
                    <el-input v-model="addform.address"></el-input>
                </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="addVisible = false">取 消</el-button>
                <el-button type="primary" @click="addEngineroom">添 加</el-button>
            </span>
        </el-dialog>

    </div>
</template>

<script>
import config from '../common/config.vue';

export default {
    data() {
        return {
            computerList:[],        // 机房数据绑定列表

            editVisible: false,     // 修改机房弹窗
            form: {},               // 修改机房绑定数据

            addVisible: false,      // 添加机房弹窗
            // 添加机房数据绑定
            addform:{
                engineroom_name: '',
                contact:'',
                contact_info:'',
                address:'',
            },

            // 添加机房、修改机房 数据验证
            rulesform:{
                engineroom_name:[{ required: true, message: '请输入机房名称', trigger: 'blur'}],
                contact:[{ required: true, message: '请输入联系人', trigger: 'blur'}],
                contact_info:[{ required: true, message: '请输入联系方式', trigger: 'blur'}],
                address:[{ required: true, message: '请输入机房地址', trigger: 'blur'}],
            }, 
       
            // 搜索
            query: {
                engineroom_name: '',
                pageIndex: 1,
                pageSize: 10,
            },
            
            // 分页
            pageTotal: 0,            
        };
    },

    methods: {
        // 添加机房
        addEngineroom() {
            this.$refs.rulesform.validate((valid) => {
                if (!valid) {
                    return false
                }
                this.$axios.post(config.urlpath+'/engineroom',this.addform).then(response => {
                    if(response.data.status == 200){
                        this.$message.success('添加成功');
                        this.searchEngineroomList();
                        this.addform = {
                            engineroom_name: '',
                            contact:'',
                            contact_info:'',
                            address:'',
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
                this.$axios.delete(config.urlpath+'/engineroom/'+row.ID).then(response => {
                    if(response.data.status == 200){
                        this.$message.success('删除成功');
                        this.computerList.splice(index, 1);
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
                this.$axios.put(config.urlpath+'/engineroom/'+this.form.ID,this.form).then(response => {
                    if(response.data.status == 200){
                        this.$message.success(`修改成功`);
                        this.editVisible = false;
                        this.searchEngineroomList()
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

        // 获取机房列表、搜索
        searchEngineroomList() {
            this.$axios.get(config.urlpath+'/engineroom', {
                params:{
                    pagesize: this.query.pageSize,
                    page: this.query.pageIndex,
                    engineroom_name: this.query.engineroom_name
                }
            }).then(response => {
                if(response.data.status == 200){
                    this.computerList = response.data.data;
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

        // 分页导航
        handlePageChange(val) {
            this.$set(this.query, 'pageIndex', val);
            this.searchEngineroomList()  
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

        // 刷新
        load() {
            this.query={
                engineroom_name: '',
                pageIndex: 1,
                pageSize: 10,
            }
            this.searchEngineroomList()
        },
    },

    // Vue生命周期函数
    created() {
        this.searchEngineroomList()
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
