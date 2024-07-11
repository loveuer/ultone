import {Component} from '@angular/core';
import {CommonModule} from '@angular/common';
import {MatFormField, MatHint, MatLabel} from "@angular/material/form-field";
import {MatOption, MatSelect} from "@angular/material/select";
import {MatInput} from "@angular/material/input";
import {MatButton} from "@angular/material/button";
import {FormsModule} from "@angular/forms";
import {NewUser} from "../../interface/user";
import {MsgService} from "../../service/msg.service";
import {
  MatDialogActions,
  MatDialogClose,
  MatDialogContent,
  MatDialogRef,
  MatDialogTitle
} from "@angular/material/dialog";
import {UserService} from "../../service/user.service";

interface user_form {
  username: string,
  nickname: string,
  password: string,
  confirm_password: string,
  status: 0 | 1,
  privileges: number[],
}

@Component({
  selector: 'app-create-user-dialog',
  standalone: true,
  imports: [
    CommonModule,
    MatFormField,
    MatSelect,
    MatOption,
    MatInput,
    MatLabel,
    MatHint,
    MatButton,
    FormsModule,
    MatDialogTitle,
    MatDialogContent,
    MatDialogActions,
    MatDialogClose,
  ],
  templateUrl: './create-user-dialog.component.html',
  styleUrl: './create-user-dialog.component.scss'
})
export class CreateUserDialogComponent {
  new_user = {
    username: "",
    nickname: "",
    password: "",
    confirm_password: "",
    status: 0,
    privileges: [],
  } as user_form

  constructor(
    public msg_srv: MsgService,
    public user_srv: UserService,
    public dialogRef: MatDialogRef<CreateUserDialogComponent>,
  ) {
  }

  create() {
    if (this.new_user.username.length < 3) {
      this.msg_srv.warning("用户名长度不能小于3")
      return
    }
    if (this.new_user.password.length < 8) {
      this.msg_srv.warning("密码至少8位")
      return
    }
    if (!this.new_user.password.match(/[a-z]+/)) {
      this.msg_srv.warning("密码必须包含小写字母")
      return
    }
    if (!this.new_user.password.match(/[A-Z]+/)) {
      this.msg_srv.warning("密码必须包含大写字母")
      return
    }
    if (this.new_user.password != this.new_user.confirm_password) {
      this.msg_srv.warning("输入的两次密码不同")
      return
    }
    console.log('[D] before create new user, new_user=', this.new_user)

    this.user_srv.manage_user_create({
      username: this.new_user.username,
      password: this.new_user.password,
      status: this.new_user.status,
      privileges: this.new_user.privileges,
      role: 100,
    } as NewUser).subscribe(val => {
      if (val.status === 200) {
        this.msg_srv.success("创建用户成功")
        this.dialogRef.close()
      }
    })
  }
}
