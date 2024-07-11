import { Injectable, Inject } from '@angular/core';
import { Component } from '@angular/core';
import { MAT_SNACK_BAR_DATA } from '@angular/material/snack-bar';
import {
  MatSnackBar,
  MatSnackBarHorizontalPosition,
  MatSnackBarVerticalPosition,
} from '@angular/material/snack-bar';
import {MatIconModule} from "@angular/material/icon";

interface SnackMessage {
  type: "success" | "error" | "warning" | "info"
  duration: number
  data: string
}

@Injectable({
  providedIn: 'root'
})
export class MsgService {
  horizontalPosition: MatSnackBarHorizontalPosition = 'center';
  verticalPosition: MatSnackBarVerticalPosition = 'top';

  constructor(
    private _snackBar: MatSnackBar
  ) { }


  success(data: string) {
    this.open({ type: "success", duration: 1500, data: data })
  }

  info(data: string) {
    this.open({ type: "info", duration: 2000, data: data })
  }

  warning(data: string) {
    this.open({ type: "warning", duration: 3000, data: data })
  }

  error(data: string) {
    this.open({ type: "error", duration: 3500, data: data })
  }

  open(msg: SnackMessage) {
    switch (msg.type) {
      case "success":
        this._snackBar.openFromComponent(SnackMessageSuccess, {
          data: msg.data,
          duration: msg.duration,
          horizontalPosition: this.horizontalPosition,
          verticalPosition: this.verticalPosition,
          panelClass: ['success-snackbar'],
        })
        break;
      case "info":
        this._snackBar.openFromComponent(SnackMessageInfo, {
          data: msg.data,
          duration: msg.duration,
          horizontalPosition: this.horizontalPosition,
          verticalPosition: this.verticalPosition,
          panelClass: ['info-snackbar'],
        })
        break
      case "warning":
        this._snackBar.openFromComponent(SnackMessageWarning, {
          data: msg.data,
          duration: msg.duration,
          horizontalPosition: this.horizontalPosition,
          verticalPosition: this.verticalPosition,
          panelClass: ['warning-snackbar'],
        })
        break
      case "error":
        this._snackBar.openFromComponent(SnackMessageError, {
          data: msg.data,
          duration: msg.duration,
          horizontalPosition: this.horizontalPosition,
          verticalPosition: this.verticalPosition,
          panelClass: ['error-snackbar'],
        })
        break
      default:
        this._snackBar.openFromComponent(SnackMessageInfo, {
          data: msg.data,
          duration: msg.duration,
          horizontalPosition: this.horizontalPosition,
          verticalPosition: this.verticalPosition,
          panelClass: ['info-snackbar'],
        })
        break;
    }
  }
}


@Component({
  selector: 'snack-message-success',
  template: `
      <div class="snack-message-success snack-message">
          <div>
              <mat-icon>cancel</mat-icon>
          </div>
          <div class="snack-message-content">
              {{ data }}
          </div>
      </div>
  `,
  standalone: true,
  imports: [
    MatIconModule
  ],
  styles: [`
    .snack-message-success {
    }
  `]
})
export class SnackMessageSuccess {
  constructor(
    @Inject(MAT_SNACK_BAR_DATA) public data: string,
    @Inject(MAT_SNACK_BAR_DATA) public duration: number,
    @Inject(MAT_SNACK_BAR_DATA) public horizontalPosition: MatSnackBarHorizontalPosition,
    @Inject(MAT_SNACK_BAR_DATA) public verticalPosition: MatSnackBarVerticalPosition,
  ) { }
}


@Component({
  selector: 'snack-message-info',
  template: `
      <div class="snack-message-info snack-message">
          <div>
              <mat-icon>cancel</mat-icon>
          </div>
          <div class="snack-message-content">
              {{ data }}
          </div>
      </div>
  `,
  standalone: true,
  imports: [
    MatIconModule
  ],
  styles: [`
    .snack-message-info {
    }
  `]
})
export class SnackMessageInfo {
  constructor(
    @Inject(MAT_SNACK_BAR_DATA) public data: string,
    @Inject(MAT_SNACK_BAR_DATA) public duration: number,
    @Inject(MAT_SNACK_BAR_DATA) public horizontalPosition: MatSnackBarHorizontalPosition,
    @Inject(MAT_SNACK_BAR_DATA) public verticalPosition: MatSnackBarVerticalPosition,
  ) { }
}

@Component({
  selector: 'snack-message-warning',
  template: `
      <div class="snack-message-warning snack-message">
          <div>
              <mat-icon>cancel</mat-icon>
          </div>
          <div class="snack-message-content">
              {{ data }}
          </div>
      </div>
  `,
  standalone: true,
  imports: [
    MatIconModule
  ],
  styles: [`
    .snack-message-warning {
    }
  `]
})
export class SnackMessageWarning {
  constructor(
    @Inject(MAT_SNACK_BAR_DATA) public data: string,
    @Inject(MAT_SNACK_BAR_DATA) public duration: number,
    @Inject(MAT_SNACK_BAR_DATA) public horizontalPosition: MatSnackBarHorizontalPosition,
    @Inject(MAT_SNACK_BAR_DATA) public verticalPosition: MatSnackBarVerticalPosition,
  ) { }
}

@Component({
  selector: 'snack-message-error',
  template: `
      <div class="snack-message-error snack-message">
          <div>
              <mat-icon>cancel</mat-icon>
          </div>
          <div class="snack-message-content">
              {{ data }}
          </div>
    </div>
  `,
  standalone: true,
  imports: [
    MatIconModule
  ],
  styles: [`
  `]
})
export class SnackMessageError {
  constructor(
    @Inject(MAT_SNACK_BAR_DATA) public data: string,
    @Inject(MAT_SNACK_BAR_DATA) public duration: number,
    @Inject(MAT_SNACK_BAR_DATA) public horizontalPosition: MatSnackBarHorizontalPosition,
    @Inject(MAT_SNACK_BAR_DATA) public verticalPosition: MatSnackBarVerticalPosition,
  ) { }
}
