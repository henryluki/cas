ckage main

import (
  "net/http"
  "io/ioutil"

  "log"
  "strings"

  "github.com/Unknwon/macaron"
  "github.com/macaron-contrib/session"
)

func main() {
  server_root := "https://nicksite.com:3000"
  local_login_url := "http://nicksite.com:4000/login"

  m := macaron.Classic()
  m.Use(macaron.Renderer())
  m.Use(session.Sessioner())

  m.Get("/login", func(ctx *macaron.Context, s session.Store) {
    ticket := ctx.Query("ticket")
    if len(ticket) == 0 {
      ctx.Redirect(server_root + "/login?service=" + local_login_url)
      return
    } else {
      s.Set("ticket", ticket)
      ctx.Redirect("/")
    }
  })

  m.Get("/", func(ctx *macaron.Context, s session.Store) string {
    if s.Get("login") != nil {
      return "Welcome, " + s.Get("login").(string)
    }

    // Retrieve the ticket
    if s.Get("ticket") == nil {
      ctx.Redirect("/login")
      return ""
    }

    // Validate the ticket
    ticket := s.Get("ticket").(string)
    resp, err := http.Get(server_root + "/validate?ticket=" + ticket + "&service=" + local_login_url)

    if err != nil {
      log.Fatalf("ERROR: %s\n", err)
      return "Unable to validate the ticket"
    }

    bs, _ := ioutil.ReadAll(resp.Body)
    split := strings.Split(string(bs), "\n")
    ticketValid, login := split[0] == "yes", split[1]

    if ticketValid {
      s.Set("login", login)
      ctx.Redirect("/")
      return ""
    }

    return "Invalid login"
  })

  m.Run()
}
