import {Component, Inject, signal} from '@angular/core';
import { CommonModule } from '@angular/common';
import {ConfirmData} from "../../interface/confirm_data";
import {MAT_DIALOG_DATA, MatDialogActions, MatDialogClose, MatDialogTitle} from "@angular/material/dialog";
import {MatButton} from "@angular/material/button";

@Component({
  selector: 'app-confirm',
  standalone: true,
  imports: [CommonModule, MatDialogTitle, MatDialogActions, MatButton, MatDialogClose],
  templateUrl: './confirm.component.html',
  styleUrl: './confirm.component.scss'
})
export class ConfirmComponent {
  constructor(@Inject(MAT_DIALOG_DATA) public data: ConfirmData) {}

  confirm() {
    this.data.confirmed = true
  }
}
