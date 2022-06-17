# GoFed

## Pkgcheck
Check which golang packages are not manageable by go-sig with the following commands:

    dnf repoquery -q  --repo=rawhide{,-source}  --whatrequires golang --recursive | grep src$ | pkgname | sort | uniq > pkgs
    cmd/pkgchk/pkgchk table -f pkgs > pkgs.table

Extract people from table

    cat pkgs.table | grep -v '+---------------+' | grep -v ' commit ' | grep -v ' admin ' | awk -F '|' '{print $3}' | sort | uniq -c | sort -n

## Email

    dnf repoquery -q  --repo=rawhide{,-source}  --whatrequires golang --recursive | grep src$ | pkgname | sort | uniq > pkgs
    cmd/pkgchk/pkgchk json -f pkgs -o pkgs.json
    cmd/tpl/tpl ml -j pkgs.json
