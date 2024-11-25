-- name: GetAllUsers :many
select *
from users
limit $1
offset $2
;

