<template>
    <div class="login-wrap">
        <div class="ms-login">
            <div class="ms-title">文件发布系统</div>
            <el-form :model="param" :rules="rules" ref="login" class="ms-content" >
                <el-form-item prop="user_name">
                    <el-input v-model="param.user_name" placeholder="用户名">
                        <el-button slot="prepend" icon="el-icon-user"></el-button>
                    </el-input>
                </el-form-item>
                <el-form-item prop="password">
                    <el-input type="password" placeholder="密码" v-model="param.password" @keyup.enter.native="submitForm()">
                        <el-button slot="prepend" icon="el-icon-lock"></el-button>
                    </el-input>
                </el-form-item>
                <div class="login-btn">
                    <el-button type="primary" @click="submitForm()">登录</el-button>
                </div>
            </el-form>
        </div>
    </div>
</template>

<script>
import config from '../common/config.vue';

export default {
    data: function() {
        return {
            param: {
                user_name: '',
                password: '',
            },
            rules: {
                user_name: [{required: true, message: '请输入用户名', trigger: 'blur'}],
                password: [{required: true, message: '请输入密码', trigger: 'blur'}],
            },
        };
    },

    methods: {
        // 登录
        submitForm() {
            this.$refs.login.validate((valid) => {
                if (valid) {
                    this.$axios.post(config.urlpath+'/login/',{
                            user_name: this.param.user_name,
                            password: this.param.password,
                    }).then(response => {
                        if(response.data.status == 200){
                            localStorage.setItem('token', response.data.token);
                            localStorage.setItem('user_name', response.data.user_name);
                            localStorage.setItem('role', response.data.role);
                            if (response.data.role == 1) {
                                this.$router.push('/');
                            } else {
                                this.$router.push('/deploymentLogs');
                            }   
                        }else{
                            this.$message({
                                type:'error',
                                message:response.data.message,
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
                } else {
                    return false;
                }
            });
        },
    },
};
</script>

<style scoped>
.login-wrap {
    position: relative;
    width: 100%;
    height: 100%;
}
.ms-title {
    width: 100%;
    line-height: 50px;
    text-align: center;
    font-size: 22px;
    color:blue;
    font-weight:bold;
    border-bottom: 1px solid #ddd;
}
.ms-login {
    position: absolute;
    left: 50%;
    top: 50%;
    width: 350px;
    margin: -190px 0 0 -175px;
    border-radius: 5px;
    background: rgba(255, 255, 255, 0.3);
    overflow: hidden;
}
.ms-content {
    padding: 30px 30px;
}
.login-btn {
    text-align: center;
}
.login-btn button {
    width: 100%;
    height: 36px;
    margin-bottom: 10px;
}
</style>