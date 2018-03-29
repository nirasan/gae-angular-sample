import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppComponent } from './app.component';
import {HttpClientModule} from "@angular/common/http";
import {CookieService} from "ngx-cookie-service";
import {AuthService} from './auth.service';
import { TodoComponent } from './todo/todo.component';
import {TodoService} from './todo.service';

@NgModule({
  declarations: [
    AppComponent,
    TodoComponent
  ],
  imports: [
    HttpClientModule,
    BrowserModule
  ],
  providers: [
    CookieService,
    AuthService,
    TodoService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
