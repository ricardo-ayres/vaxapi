# VaxApi

Projeto Integrador do segundo semestre de 2022, Eixo da Computação, Univesp.

PJI310 - Sala 001 - Grupo 011

## O projeto:

O conceito que explorado neste projeto foi o de uma plataforma de carteiras de
vacinação virtuais, em que o banco de dados armazena os dados dos usuários,
das vacinas e com quais doses de quais vacinas cada usuário foi vacinado.

## O Modelo de dados:

O modelo de dados para este projeto foi pensado com base nas seguintes
necessidades: - Cada vacina precisa ter um cadastro que determine o nome
da vacina, quantas doses são necessárias e qualquer outra informação
particular daquela vacina; - Cada usuário precisa ter um cadastro que
mantenha suas informações básicas de identificação; - Há necessidade
de manter um registro de cada dose de cada vacina que cada um dos usuários
já tomou e a data em que aquela dose foi aplicada.

Com base nessas informações estabeleceu-se a necessidade de uma relação
de muitos para muitos: Cada usuário se relaciona com múltiplas vacinas
(mais de uma vez, inclusive), e existem vários usuários. Com esta conclusão
montou-se um banco de dados com 3 tabelas:

- Tabela de usuários: Armazena as informações pessoais básicas e de
identificação e acesso de cada usuário;

- Tabela de vacinas: Armazena as informações de cada vacina, como nome
da vacina, número de doses e outras informações relevantes;

- Tabela de doses: Armazena por meio de chaves estrangeiras (*foreign
keys*) as doses de quais vacinas tomadas por cada usuário e as datas em
que aquelas doses foram aplicadas.

Desta forma obteve-se um banco de dados suficientemente normalizado para
o que se propõe a fazer neste projeto, e que preserva o requerimento de
relações do tipo muitos para muitos.

## O API:

Nosso projeto é uma API REST escrita utilizando os módulos da biblioteca
padrão da linguagem Go e banco de dados SQLite3, utilizando o driver
[go-sqlite](https://www.github.com/glebarez/go-sqlite).

A nossa API consiste em implementações simples das operações CRUD no
banco de dados, expostas na forma de uma API REST.  Um usuário deve ser
capaz de criar seu próprio usuário, modificálo e excluí-lo se assim
desejar. Estas operações são expostas da seguinte maneira:

A criação de novos usuários é feita enviando uma requisição `POST`
para o caminho `/users/` e contendo como corpo os dados do usuário em
formato json:

```
{
	"username": "nome_do_novo_usuário",
	"name": "Nome completo do novo usuário",
	"birth": "aaaa-mm-dd (data de nascimento)",
	"email": "email@do_usuário.com",
	"password": "senha do novo usuário a ser criado"
}
```

Caso a requisição obtenha êxito e o usuário seja criado o sistema irá
responder com o status 200 "OK" e irá enviar no corpo da resposta uma
cópia dos dados do usuário recém criado, em formato json e excluindo a
senha e contendo o id de usuário.


Para consulta de todos os tipos de vacinas registradas no sistema consulta-se
`/vacs/` com uma requisição `GET`.

Todas as requisições a seguir necessitam de autenticação tipo `Basic`
no cabeçalho da requisição.

Para consulta dos dados do usuário, faz-se uma requisição `GET` em
`/users/`.  Em caso de sucesso a resposta será como no caso de criação de
novos usuários, uma cópia dos dados do usuário será enviada em formato
json. O usuário requisitado será definido com base nas informações de
autenticação no cabeçalho.

A alteração dos dados do usuário é feita acessando-se o recurso `/users/`
com o método `PUT` e o corpo da mensagem deve conter as informações do
usuário a serem atualizadas, em formato json:

```
{
	"email": "troquei@de_email.com"
}
```

Em caso de sucesso a resposta será como no caso de criação de novos
usuários, uma cópia dos dados do usuário será enviada em formato json.

Para remoção de um usuário acessa-se o recurso `/users/` com o método
`DELETE`. O usuário requisitado será definido com base nas informações
de autenticação no cabeçalho. O servidor responderá com 200 "OK" em
caso de sucesso e o corpo da mensagem será vazio.

Para consulta de todas as doses de todas as vacinas que o usuário já tomou,
acessa-se com `GET` o recurso `/users/doses/`.

Para consulta de todas as doses de uma vacina em específico que um dado
usuário já tomou, acessa-se com `GET` o recurso `/users/doses/[vid]`
em que `[vid]` é o identificador numérico da vacina que se deseja consultar.

Para registro de uma nova dose, acessa-se `/users/doses/` com o método
`POST` e o corpo da mensagem contendo o id numérico da vacina e a data em
que a dose foi tomada, em formato json:

```
{
	"vac_id": [vid],
	"date_taken": "aaaa-mm-dd"
}
```

Em que `[vid]` é o identificador numérico da vacina. Em caso de sucesso o
sistema retornará 200 "OK" com uma cópia das informações da dose recém
registrada em formato json.

Para alteração ou remoção de uma dose utiliza-se o mesmo protocolo,
porém com os
métodos `PUT` ou `DELETE` e especificando-se no corpo da mensagem o
identificador numérico `[did]` da dose tomada, em formato json:

```
{
	"dose_id": [did]
}
```

O cadastro, alteração e remoção de vacinas não foi implementado, visto
que estas informações devem ser controladas apenas pela administração
e foge ao escopo do projeto um sistema de administração do banco de dados.
