import {Component, OnInit} from '@angular/core';
import {Todo} from '../todo';
import {TodoService} from '../todo.service';


@Component({
  selector: 'app-todo',
  templateUrl: './todo.component.html',
  styleUrls: ['./todo.component.css']
})
export class TodoComponent implements OnInit {

  todoList: Todo[] = [];

  constructor(private todo: TodoService) {
  }

  ngOnInit() {
    this.getList();
  }

  getList() {
    this.todo.getList().subscribe(data => this.todoList = data || []);
  }

  create(content: string) {
    const t = {
      id: 0,
      user_id: '',
      done: false,
      content: content
    };
    this.todo.create(t).subscribe(data => this.todoList.push(data));
  }

  delete(t: Todo) {
    this.todoList = this.todoList.filter(tt => tt !== t);
    console.log(this.todoList);
    this.todo.delete(t).subscribe();
  }

}
