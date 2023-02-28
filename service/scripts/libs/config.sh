#! /bin/bash -e

# Namespace Constant
BLOGRPC_NAMESPACE='blogrpc'
SERVICE_NAMESPACE="${BLOGRPC_NAMESPACE}/service"

##################################################
# Get root path
# Arguments:
#   None
# Returns:
#   String
##################################################
get_root_path() {
  local current_path=$(pwd)
  cd ${0%/*} # cd to ~/blogrpc/scripts/
  cd ../../
  local root_path=$(pwd)
  cd ${current_path}
  echo ${root_path}
}

###################################################
# Get service list
# Arguments:
#   None
# Returns:
#   String
###################################################
get_service_list() {
  declare -a services
  declare -a blogrpc_files

  # get all files
  local blogrpc_path="$(get_root_path)/service"
  blogrpc_files=(${blogrpc_path}/*)

  # get service name
  for file in "${blogrpc_files[@]}"; do
    # if a folder contains main.go, then it is a service folder
    if [[ -f "$file/main.go" ]]; then
      services+=(${file##*/})
    fi
  done

  echo ${services[@]}
}

ROOT_PATH="$(get_root_path)"
PROTO_PATH="${ROOT_PATH}/proto"
# GOPATH_IN_DOCKER is the GOPATH value in docker golang builder container
GOPATH_IN_DOCKER='/app'

# SERVICE_LIST is the RPC repos to be pulled to build
SERVICE_LIST=($(get_service_list))