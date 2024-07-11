import {Component, OnInit} from '@angular/core';
import {CommonModule} from '@angular/common';
import {UserService} from "../../service/user.service";
import {MatButton, MatIconButton} from "@angular/material/button";
import {MatDialog} from "@angular/material/dialog";
import {CreateUserDialogComponent} from "../../component/create-user-dialog/create-user-dialog.component";
import {
  MatCell,
  MatCellDef,
  MatColumnDef,
  MatHeaderCell,
  MatHeaderCellDef,
  MatHeaderRow, MatHeaderRowDef, MatRow, MatRowDef,
  MatTable
} from "@angular/material/table";
import {MatIcon} from "@angular/material/icon";
import {Enum} from "../../interface/enum";
import {User} from "../../interface/user";
import {ConfirmComponent} from "../../component/confirm/confirm.component";
import {UpdateUserDialogComponent} from "../../component/update-user-dialog/update-user-dialog.component";
import {MatPaginator, PageEvent} from "@angular/material/paginator";

@Component({
  selector: 'app-user',
  standalone: true,
  imports: [
    CommonModule,
    MatButton,
    MatTable,
    MatColumnDef,
    MatHeaderCell,
    MatCell,
    MatCellDef,
    MatHeaderCellDef,
    MatHeaderRow,
    MatRow,
    MatRowDef,
    MatHeaderRowDef,
    MatIconButton,
    MatIcon,
    MatPaginator,
  ],
  templateUrl: './user.component.html',
  styleUrl: './user.component.scss'
})
export class UserComponent implements OnInit {
  readonly displayedColumns: string[] = ["username", "status", "role", "privileges", "operation",];

  constructor(
    public user_srv: UserService,
    public dialog: MatDialog,
  ) {
  }

  ngOnInit() {
    this.user_srv.manage_user_list()
  }

  open_dialog() {
    this.dialog.open(CreateUserDialogComponent, {data: {}})
  }

  _parsePrivileges(privileges: any): string {
    try {
      let ps = privileges as Enum<number>[]
      return ps.map(v => v.label).join("; ")
    } catch (e) {
      console.log("[D] parse privileges err=", e)
      return ""
    }
  }

  delete_user(target: User) {
    let data = {title: `确认删除用户 ${target.username} 吗?`, confirmed: false }
    this.dialog.open(ConfirmComponent, {data: data}).afterClosed().subscribe(() => {
      if (data.confirmed) {
        this.user_srv.manage_user_delete(target)
      }
    })
  }

  update_user(target: User) {
    this.dialog.open(UpdateUserDialogComponent, {data: target})
  }

  handlePager(event: PageEvent) {
    console.log('[D] handle pager change event=', event)
    this.user_srv.user_list.set({...this.user_srv.user_list(), size: event.pageSize, page: event.pageIndex})
    this.user_srv.manage_user_list()
  }
}

