import {Component, Inject} from '@angular/core';
import {CommonModule} from '@angular/common';
import {MatButton} from "@angular/material/button";
import {
  MAT_DIALOG_DATA,
  MatDialogActions,
  MatDialogClose,
  MatDialogContent, MatDialogRef,
  MatDialogTitle
} from "@angular/material/dialog";
import {MatFormField, MatLabel} from "@angular/material/form-field";
import {MatInput} from "@angular/material/input";
import {MatOption} from "@angular/material/autocomplete";
import {MatSelect} from "@angular/material/select";
import {FormControl, FormsModule, ReactiveFormsModule} from "@angular/forms";
import {User} from "../../interface/user";
import {Enum} from "../../interface/enum";
import {MsgService} from "../../service/msg.service";
import {ArrayEquals} from "../../../tool";
import {UserService} from "../../service/user.service";

interface updateUser extends User {
  password: string
  confirm_password: string
}

type Updates = { [key: string]: any }

@Component({
  selector: 'app-update-user-dialog',
  standalone: true,
  imports: [CommonModule, MatButton, MatDialogActions, MatDialogClose, MatDialogContent, MatDialogTitle, MatFormField, MatInput, MatLabel, MatOption, MatSelect, ReactiveFormsModule, FormsModule],
  templateUrl: './update-user-dialog.component.html',
  styleUrl: './update-user-dialog.component.scss'
})
export class UpdateUserDialogComponent {

  user: updateUser = {} as updateUser

  privileges = new FormControl([] as number[])
  status = new FormControl()
  password = new FormControl('')
  confirm_password = new FormControl('')

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: User,
    private msg_srv: MsgService,
    public user_srv: UserService,
    public dialogRef: MatDialogRef<UpdateUserDialogComponent>,
  ) {
    this.user = {...data, password: '', confirm_password: ''}
    this.privileges.setValue(this.user.privileges.map(v => v.value))
    this.status.setValue(this.user.status.value)
  }

  update() {
    let updated = false
    let updates: Updates = {"id": this.user.id}
    if (this.password.value || this.confirm_password.value) {
      updated = true
      updates["password"] = this.password.value
      if (this.password.value != this.confirm_password.value) {
        this.msg_srv.warning('两次密码不相同')
        return
      }
    }

    if (this.status.value != this.user.status.value) {
      updated = true
      updates["status"] = this.status.value
    }

    if (!ArrayEquals(this.privileges.value!, this.user.privileges.map(p => p.value))) {
      updates["privileges"] = this.privileges.value
      updated = true
    }

    if (!updated) {
      this.msg_srv.warning('没有变更')
      return
    }

    this.user_srv.manage_user_update(updates).subscribe(rs => {
      if (rs.status === 200 || rs.status === 401) {
        this.dialogRef.close()
      }
    })
  }

  initedPrivileges(privileges: Enum<number>[]) {
    return privileges.map(p => p.value)
  }
}
