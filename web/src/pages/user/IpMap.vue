<template>
    <div>
        <el-card>

            <el-form :inline="true">
                <el-form-item>
                    <el-button
                            size="small"
                            type="primary"
                            icon="el-icon-plus"
                            @click="handleEdit('')">添加IP保留
                    </el-button>
                </el-form-item>
                <el-form-item>
                    <el-alert
                            title="IP保留功能：为指定用户名永久保留IP地址，该用户每次连接都会获得相同的IP"
                            type="info"
                            :closable="false">
                    </el-alert>
                </el-form-item>
            </el-form>

            <el-table
                    ref="multipleTable"
                    :data="tableData"
                    border>

                <el-table-column
                        sortable="true"
                        prop="id"
                        label="ID"
                        width="60">
                </el-table-column>

                <el-table-column
                        prop="ip_addr"
                        label="IP地址"
                        width="150">
                </el-table-column>

                <el-table-column
                        prop="username"
                        label="用户名"
                        width="150">
                </el-table-column>

                <el-table-column
                        prop="keep"
                        label="IP保留"
                        width="100">
                    <template slot-scope="scope">
                        <el-switch
                                disabled
                                v-model="scope.row.keep"
                                active-color="#13ce66">
                        </el-switch>
                    </template>
                </el-table-column>

                <el-table-column
                        prop="mac_addr"
                        label="MAC地址"
                        width="150">
                </el-table-column>

                <el-table-column
                        prop="unique_mac"
                        label="唯一MAC"
                        width="100">
                    <template slot-scope="scope">
                        <el-tag v-if="scope.row.unique_mac" type="success">是</el-tag>
                        <el-tag v-else type="info">否</el-tag>
                    </template>
                </el-table-column>

                <el-table-column
                        prop="note"
                        label="备注">
                </el-table-column>

                <el-table-column
                        prop="last_login"
                        label="最后登陆时间"
                        :formatter="tableDateFormat">
                </el-table-column>

                <el-table-column
                        label="操作"
                        width="150">
                    <template slot-scope="scope">
                        <el-button
                                size="mini"
                                type="primary"
                                @click="handleEdit(scope.row)">编辑
                        </el-button>

                        <el-popconfirm
                                class="m-left-10"
                                @confirm="handleDel(scope.row)"
                                title="确定要删除IP映射吗？">
                            <el-button
                                    slot="reference"
                                    size="mini"
                                    type="danger">删除
                            </el-button>
                        </el-popconfirm>

                    </template>
                </el-table-column>
            </el-table>

            <div class="sh-20"></div>

            <el-pagination
                    background
                    layout="prev, pager, next"
                    :pager-count="11"
                    @current-change="pageChange"
                    :total="count">
            </el-pagination>

        </el-card>

        <!--新增、修改弹出框-->
        <el-dialog
                title="IP映射管理"
                :close-on-click-modal="false"
                :visible="user_edit_dialog"
                @close="disVisible"
                width="600px"
                center>

            <el-form :model="ruleForm" :rules="rules" ref="ruleForm" label-width="120px" class="ruleForm">
                <el-form-item label="ID" prop="id">
                    <el-input v-model="ruleForm.id" disabled></el-input>
                </el-form-item>
                <el-form-item label="IP地址" prop="ip_addr">
                    <el-input v-model="ruleForm.ip_addr"></el-input>
                </el-form-item>
                <el-form-item label="用户名" prop="username">
                    <el-input v-model="ruleForm.username"></el-input>
                </el-form-item>
                
                <el-form-item label="IP保留" prop="keep">
                    <el-switch
                            v-model="ruleForm.keep"
                            active-color="#13ce66">
                    </el-switch>
                    <span style="margin-left: 10px; color: #909399; font-size: 12px;">
                        开启后，该IP将永久分配给指定用户名
                    </span>
                </el-form-item>

                <el-form-item label="MAC地址" prop="mac_addr">
                    <el-input v-model="ruleForm.mac_addr" placeholder="自动记录，无需手动填写"></el-input>
                </el-form-item>

                <el-form-item label="备注" prop="note">
                    <el-input v-model="ruleForm.note"></el-input>
                </el-form-item>

                <el-form-item>
                    <el-button type="primary" @click="submitForm('ruleForm')">保存</el-button>
                    <el-button @click="disVisible">取消</el-button>
                </el-form-item>
            </el-form>

        </el-dialog>

    </div>
</template>

<script>
import axios from "axios";

export default {
    name: "IpMap",
    components: {},
    mixins: [],
    created() {
        this.$emit('update:route_path', this.$route.path)
        this.$emit('update:route_name', ['用户信息', 'IP映射'])
    },
    mounted() {
        this.getData(1)
    },
    data() {
        return {
            tableData: [],
            count: 10,
            nowIndex: 0,
            ruleForm: {
                status: 1,
                groups: [],
            },
            rules: {
                username: [
                    {required: true, message: '请输入用户名', trigger: 'blur'},
                    {max: 60, message: '长度小于 60 个字符', trigger: 'blur'}
                ],
                mac_addr: [
                    {required: false, message: 'MAC地址自动记录', trigger: 'blur'}
                ],
                ip_addr: [
                    {required: true, message: '请输入IP地址', trigger: 'blur'}
                ],
                status: [
                    {required: true}
                ],
            },
        }
    },
    methods: {
        getData(p) {
            axios.get('/user/ip_map/list', {
                params: {
                    page: p,
                }
            }).then(resp => {
                var data = resp.data.data
                console.log(data);
                this.tableData = data.datas;
                this.count = data.count
            }).catch(error => {
                if (error.response && error.response.status === 401) {
                    return;
                }
                this.$message.error('哦，请求出错');
                console.log(error);
            });
        },
        pageChange(p) {
            this.getData(p)
        },
        handleEdit(row) {
            !this.$refs['ruleForm'] || this.$refs['ruleForm'].resetFields();
            console.log(row)
            this.user_edit_dialog = true
            if (!row) {
                return;
            }

            axios.get('/user/ip_map/detail', {
                params: {
                    id: row.id,
                }
            }).then(resp => {
                this.ruleForm = resp.data.data
            }).catch(error => {
                if (error.response && error.response.status === 401) {
                    return;
                }
                this.$message.error('哦，请求出错');
                console.log(error);
            });
        },
        handleDel(row) {
            axios.post('/user/ip_map/del?id=' + row.id).then(resp => {
                var rdata = resp.data
                if (rdata.code === 0) {
                    this.$message.success(rdata.msg);
                    this.getData(1);
                } else {
                    this.$message.error(rdata.msg);
                }
                console.log(rdata);
            }).catch(error => {
                if (error.response && error.response.status === 401) {
                    return;
                }
                this.$message.error('哦，请求出错');
                console.log(error);
            });
        },
        submitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (!valid) {
                    console.log('error submit!!');
                    return false;
                }

                // alert('submit!');
                axios.post('/user/ip_map/set', this.ruleForm).then(resp => {
                    var rdata = resp.data
                    if (rdata.code === 0) {
                        this.$message.success(rdata.msg);
                        this.getData(1);
                    } else {
                        this.$message.error(rdata.msg);
                    }
                    console.log(rdata);
                }).catch(error => {
                    this.$message.error('哦，请求出错');
                    console.log(error);
                });
            });
        },
    },
}
</script>

<style scoped>

</style>