# AVIA

A - 
V
I
A

The script for easy create and update new branches with base the main branch.

# Minimum requirements 
- [GIT](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- [jq](https://stedolan.github.io/jq/download/)
- [curl](https://curl.haxx.se)

## Install


copy 'avia' the script in your directory '/bin', above open the terminal(shell) and
call for ```$ avia```, done!

### Copy the envoriments

```
export TOKEN=<token-open-project>
export BRANCH_UPLEVEL=<branch-default-development>
export SERVER_OP=<ip-open-project>
export PROJECT_ID=<id-project>
```

## commands

For more informations

```
$ avia --help
```

# Timer

quando tiver começado a atividade, pode-se adicionar pausa ou continuar o tempo da mesma clicando 'p' para pausar e 'c' para constinuar. caso tenha finalizado só precionar 'd', assim ele fecha a branch e adiciona o tempo no open project

## TO-DO

- adicionar horas trabalhadas
  - [x] dicionar tbm pausa, tipo indo para o banheiro.. poderia ver tbm quando esta em lock ele colocaria pausa e se saisse ele tiraria a pausa
  - [x] quando finalizar a task ele salva no OP as informações das horas trabalhadas e tals
  - [x] fazer um facilitador para instalar e adicionar as informações como token e nome do projeto, podemos fazer no .zshrc/.bashrc e o carinha pega dela as informações de token e tals
  - criar fluxograma de como ele deve funcionar para mostrar a professora e presentar no meetup
