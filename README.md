# golang-project
Golang + Elasticsearch + Mysql

- Go Lang
- Rabbitmq 3.7.12
- Erlang 21.2
- Mysql
- Elastic V3


Mysql query

CREATE TABLE `news` (
  `id` int(20) NOT NULL,
  `author` text,
  `body` text,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
