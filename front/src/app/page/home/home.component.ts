import {Component} from '@angular/core';
import {CommonModule} from '@angular/common';
import {UserService} from "../../service/user.service";

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './home.component.html',
  styleUrl: './home.component.scss'
})
export class HomeComponent {
  constructor(
    public user_srv: UserService
  ) {
  }
}
