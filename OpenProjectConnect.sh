#!/bin/bash

# carregando variaveis de configurações
loadFileConfig() {
  list=$(awk '/=/ {print $1}' credentials.op.conf)

  # My input source is the contents of a variable called $list #
  while IFS= read -r pkg
  do
    VALUE=${pkg#*=}
    # printf 'Installing php package %s...\n' "${pkg%=$VALUE} $VALUE"
    eval "${pkg%=$VALUE}"="$VALUE"
    # /usr/bin/apt-get -qq install $pkg
  done <<< "$list"
}

# conectando a api do open project e pegando o nome da task
getNameTaskOpenProject() {
  TASK_NAME=$(curl -u apikey:$TOKEN http://$SERVER_OP/api/v3/work_packages/$1 | jq '.subject' )
  
  # formatando o nome para apresentar
  TASK_NAME=${TASK_NAME//  / }
  TASK_NAME=${TASK_NAME//\\/}
  TASK_NAME=${TASK_NAME//\"/}
  TASK_NAME=${TASK_NAME// /-}
}
