# splicepipe

sender -> splicepipe -> receiver

## usage

pass in the input and output fifo paths.
the fifos will be creates for you.

```sh
splicepipe /path/to/input/fifo /path/to/output/fifo
```

## why

i wanted to plot some data in a tool.

the [plotting tool i'm using](https://github.com/cactusdynamics/wesplot) renders the plots in a web page.
this is nice, since it lets me open the plots on my laptop.
however restarting the plotter requires refreshing the page.

this can be fixed with a fifo.
i can point my tool at the fifo to write, read the fifo into the plotter.

```sh
# make fifo
mkfifo /tmp/fifo

# start the plotter, forking to background
wesplot < /tmp/fifo &

mytool > /tmp/fifo
```

however when the peer of the fifo closes the file, it closes for both sides.
directing the data into the plotter this way means it wont be reopened, so the plotter needs to be restarted.

directing an empty descriptor can prevent the reader from closing.
as long as there is at least one writer, the fifo will remain open.

```sh
3>/tmp/fifo &
```

to get around this we can write a tool that reads the fifo, piping the data to the plotter, and reopening the fifo whenever it closes.

```sh
fiforeader /tmp/fifo | wesplot
```

we have the same problem in the other direction. when the reader closes the fifo, the writer is closed too.
having to restart the tool any time i want to make changes to the plotter is also undesirable.

to fix this we can make a similar `fifowriter`.
we could also make a fifo coupling.
two separate fifos, read and written to by a middle program, which continues running and reopens the fifo handles any time a writer or reader closes.

```sh
splicepipe /tmp/input /tmp/output

wesplot < /tmp/output &

mytool > /tmp/intput
```

---

i also needed to add a domain column to the data.
due to how the plotter functions, i cannot add it to my tool, since the tool would reset to zero each time.
this would produce bad plots.
to fix this, the domain, which is just the line number, is added in the coupler.
this way when the writer restarts, the plotter doesn't see that they started over.
