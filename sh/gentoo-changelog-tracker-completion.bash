_gentoo_changelog_tracker_complete() {
  export EIX_LIMIT=0
  local cur prev opts packages
  COMPREPLY=()
  cur="${COMP_WORDS[COMP_CWORD]}"
  prev="${COMP_WORDS[COMP_CWORD-1]}"

  if [[ "$cur" == */* ]]; then
      cat="${cur%%/*}"
      atom="${cur#*/}"
  else
      atom=""
      cat="$cur"
  fi

  echo "${cat} / ${atom} ==" >> /tmp/toto
  packages=$(eix -#n -C "$cat" "${atom}" | awk '{print $1}' | sort)

  if [[ ${prev} == "gentoo-changelog-tracker" ]]; then
    COMPREPLY=( $(compgen -W "${packages}") )
  fi
  return 0
}

complete -F _gentoo_changelog_tracker_complete gentoo-changelog-tracker
