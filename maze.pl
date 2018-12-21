
location(twisty_little, tl).
location(little_twisty, lt).
location(twenty, tw).
location(dtwisty, dt).
location(thirsty, th).
location(cabbages, cab).
location(weight_room, wr).
location(health_and_fitness, hf).


path(tl, nw, lt).
path(lt, up, th).
path(th, d, lt).
path(th, se, tw).
path(tw, w, lt).
path(tw, d, dt).
path(dt, nw, cab).
path(cab, n, lt).
path(cab, nw, hf).
path(cab, ne, wr).



navigate(From, To, Path) :- navigate(From, To, Path, []).

navigate(A, A, [], _).

navigate(From, To, [Step | Rest], Visited) :-
    path(From, Step, Intermediate),
    not(member(Intermediate, Visited)),
    navigate(Intermediate, To, Rest, [Intermediate | Visited]).


connected(A, B) :- navigate(A, B, _).