import {Component, OnInit} from '@angular/core';
import {CommonModule} from '@angular/common';
import {LogService} from "../../service/log.service";
import {decode} from "js-base64"
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
import {MatIconButton} from "@angular/material/button";

@Component({
  selector: 'app-log',
  standalone: true,
  imports: [CommonModule, MatTable, MatColumnDef, MatHeaderCell, MatCell, MatCellDef, MatHeaderCellDef, MatIcon, MatIconButton, MatHeaderRow, MatHeaderRowDef, MatRowDef, MatRow],
  templateUrl: './log.component.html',
  styleUrl: './log.component.scss'
})
export class LogComponent implements OnInit {
  displayedColumns = ["username", "created_at", "type", "html"]

  constructor(
    public log_srv: LogService,
  ) {
  }

  ngOnInit() {
    this.log_srv.get_logs()
  }

  decodeLog(org: string): string {
    return decode(org)
    //.from(org, "base64").toString()
  }
}
