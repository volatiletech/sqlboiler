SELECT * FROM "a" WHERE (a=$1 or b=$2) AND (c=$3) GROUP BY id, name HAVING id <> $4 AND length(name, $5) > $6;
