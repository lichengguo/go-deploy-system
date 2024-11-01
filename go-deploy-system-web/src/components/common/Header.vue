<template>
    <div class="header">
        <!-- 折叠按钮 -->
        <div class="collapse-btn" @click="collapseChage">
            <i v-if="!collapse" class="el-icon-s-fold"></i>
            <i v-else class="el-icon-s-unfold"></i>
        </div>
        <div class="logo">文件发布系统</div>
        <div class="header-right">
            <div class="header-user-con">
                <!-- 用户头像 -->
                <div class="user-avator">
                    <img v-bind:src="head" />
                </div>
                <!-- 用户名下拉菜单 -->
                <el-dropdown class="user-name" trigger="click" @command="handleCommand">
                    <span class="el-dropdown-link">
                        {{username}}
                        <i class="el-icon-caret-bottom"></i>
                    </span>
                    <el-dropdown-menu slot="dropdown">
                        <el-dropdown-item divided command="loginout" >退出登录</el-dropdown-item>
                         <el-dropdown-item divided command="changePassword">修改密码</el-dropdown-item>
                    </el-dropdown-menu>
                </el-dropdown>
            </div>
        </div>
        
        <!-- 修改密码 -->
        <el-dialog title="修改密码" :visible.sync="centerDialogVisible" width="30%" center>
            <el-form :model="ruleForm" :rules="rules" ref="ruleForm" hide-required-asterisk>
                <el-form-item label="旧密码" prop="old_pwd">
                    <el-input v-model="ruleForm.old_pwd" type="password"></el-input>
                </el-form-item>
                <el-form-item label="新密码" prop="new_pwd">
                    <el-input v-model="ruleForm.new_pwd" type="password"></el-input>
                </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="unChangePassword('ruleForm')">取 消</el-button>
                <el-button type="primary" @click="changePassword('ruleForm')">确 定</el-button>
            </span>
        </el-dialog>
    </div>
</template>
<script>
import bus from '../common/bus';
import config from '../common/config.vue';

export default {
    data() {
        return {
            centerDialogVisible: false,
            collapse: false,
            fullscreen: false,
            name: '',
            head:!localStorage.getItem('head') ? require("../../assets/img/img.jpg") :config.urlpath + "/" +localStorage.getItem('head'),
            ruleForm: {
                old_pwd:'',
                new_pwd:'',
            },
            // 修改密码校验
            rules: {
                old_pwd: [
                    { required: true, message: '请输入旧密码', trigger: 'blur' },
                    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
                ],
                new_pwd: [
                    { required: true, message: '请输入新密码', trigger: 'blur' },
                    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
                ]
            },
        };
    },
    


    computed: {
        username() {
            let username = localStorage.getItem('user_name');
            return username ? username : this.name;
        }
    },
    methods: {
        // 用户名下拉菜单选择事件
        handleCommand(command) {
            var that = this;
            if (command == 'loginout') {
                that.$message.success('安全退出');
                localStorage.removeItem('token');
                that.$router.push('/login');
            }
            if (command == 'changePassword' ) {
                that.centerDialogVisible = true
            }
            
        },
        // 侧边栏折叠
        collapseChage() {
            this.collapse = !this.collapse;
            bus.$emit('collapse', this.collapse);
        },

        // 修改密码
        changePassword(formName){
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    // 发送请求给服务端
                    this.$axios.post(config.urlpath+'/user/changepassword', {
                        old_pwd:this.ruleForm.old_pwd, 
                        new_pwd:this.ruleForm.new_pwd
                    }).then(response => {
                        if (response.data.status == 200) {
                            localStorage.removeItem('token', response.data.token);
                            this.$router.push('/login');
                            this.$message({
                                message:"密码修改成功, 请重新登录",
                                type:"success",
                            })
                        } else {
                            this.$message({
                                message:response.data.message+', 密码修改失败',
                                type:"error",
                                showClose: true, 
                                duration: 0,
                            })
                            return
                        }
                    }).catch(err => {
                        console.log(err)
                        alert("密码修改失败")
                    })
                    this.centerDialogVisible = false
                    this.unChangePassword()
                } else {
                    alert('校验不通过')
                }
            }) 
        },

        // 取消修改密码
        unChangePassword(){
            this.ruleForm.old_pwd = ""
            this.ruleForm.new_pwd = ""
            this.centerDialogVisible = false
        }
        
    },
    mounted() {
        if (document.body.clientWidth < 1500) {
            this.collapseChage();
        }
    }
};
</script>
<style scoped>
.header {
    position: relative;
    box-sizing: border-box;
    width: 100%;
    height: 60px;
    font-size: 18px;
    color: #fff;
}
.collapse-btn {
    float: left;
    padding: 0 23px;
    cursor: pointer;
    line-height: 60px;
}
.header .logo {
    float: left;
    width: 250px;
    line-height: 60px;
}
.header-right {
    float: right;
    padding-right: 50px;
}
.header-user-con {
    display: flex;
    height: 60px;
    align-items: center;
}
.user-name {
    margin-left: 10px;
}
.user-avator {
    margin-left: 20px;
}
.user-avator img {
    display: block;
    width: 40px;
    height: 40px;
    border-radius: 50%;
}
.el-dropdown-link {
    color: #fff;
    cursor: pointer;
}
</style>
