let homePage = 1;
let Fetching = true;
let scrollEnd = false;

async function display_post(page, api, urlParams) {
  let like = urlParams?.like || "";
  let date = urlParams?.date || "";
  let category = urlParams?.category || "";

  const { state, _ } = await window.checkIfLoggedIn();

  if ((like != "" || date != "") && !state) {
    // window.location.replace("/error?error=404")
    return;
  } else if (
    (date != "" && date != "ASC" && date != "DESC") ||
    (like != "" && like != "ASC" && like != "DESC")
  ) {
    // window.location.replace("/error?error=404")
    return;
  }

  const element = document.getElementById("centerHome");

  try {
    const response = await fetch(
      `/api/${api}?page=${page}&categ=${category}&date=${date}&like=${like}`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      }
    );
    if (!response.ok) {
      dangerError(response.status);
      return;
    }

    let postsData = await response.json();

    if (postsData !== null) {
      for (const postData of postsData) {
        const date = new Date(postData.created_at);
        const options = { year: "numeric", month: "long", day: "numeric" };
        const formattedDate = date.toLocaleDateString("en-US", options);
        let optLike = "";
        let optDisLike = "";
        if (state) {
          try {
            const respPostLike = await fetch(
              `/api/checklike?post_id=${postData.id}`
            );
            if (respPostLike.ok) {
              const respPostLikeJson = await respPostLike.json();
              if (respPostLikeJson !== null) {
                if (respPostLikeJson.like === 1) {
                  optLike = "liked";
                } else {
                  optDisLike = "disliked";
                }
              }
            }
          } catch (error) {
            console.error("Fetch failed:", error);
            createAlert(alert("alert-danger", "feald to display post"));
          }
        }
        let x = `
                <div class="myPosts">
                    <div class="postProfile" id="pid${postData.id}">
                            <div class="postImg">
                            <a><img src="/public/img/default.avif" alt=""></a>
                            </div>
                            <div class="usernameProfil">
                                <div>${postData.username}</div>
                                <div class="postDate">${formattedDate}&nbsp;&nbsp;<i class="fa-solid fa-boxes-stacked"></i>
                                    ${
                                      postData.categories
                                        ? postData.categories
                                        : "General"
                                    }
                                </div>
                            </div>
                        </div>
                        <div class="postTitle">
                            <h2>${postData.title ? postData.title : ""}</h2>
                        </div>
                        <div class="postContent">
                                ${
                                  postData.content
                                    ? postData.content.replace(/\n/g, "</br>")
                                    : ""
                                }
                        </div>
                        ${
                          postData.image
                            ? '<div class="PostImage"><img src="' +
                              postData.image +
                              '" alt=""></div>'
                            : ""
                        }
                        <div class="reaction">
                            <div>
                            <i id="pid${postData.id}l" onclick="postLike(${
          postData.id
        },'liked')" class="fa-regular fa-thumbs-up cursor ${optLike}"></i>&nbsp;<span id="pid${
          postData.id
        }ls" >${postData.likes}</span>
                            </div>
                            <div>
                            <!-- use "liked" "disliked" class -->
                            <i id="pid${postData.id}d" onclick="postLike(${
          postData.id
        },'disliked')" class="fa-regular fa-thumbs-down cursor ${optDisLike}"></i>&nbsp;<span id="pid${
          postData.id
        }ds">${postData.dislikes}</span>
                            </div>
                            <div>
                            |
                            </div>
                            <div class="cursor" onclick="toggleComment('${
                              postData.id
                            }')">
                                <i class="fa-regular fa-comment"></i>&nbsp;<span id="pid${
                                  postData.id
                                }cs">${postData.comments}</span>
                            </div>
                        </div>
                        <div id="c${postData.id}" class="none showComments">
                            <div id="c${postData.id}l" class="comment">
                                <div class="guestAddComment">
                                    <h3>Log in or sign up to share your thoughts!</h3>
                                    <a class="btn btnNotify" href="/log-in">Log in</a>
                                </div>
                            </div>

                            <div id="c${
                              postData.id
                            }a" class="comment addComment">
                                <form class="addComForm" method="post" id="c${
                                  postData.id
                                }af">
                                    <div class="postProfile">
                                        <!-- code added here in js -->
                                    </div>
                                    <div class="addpostTitle">
                                        <input type="hidden" name="post_id" value="${
                                          postData.id
                                        }">
                                        <textarea  id="c${
                                          postData.id
                                        }con" maxlength="1000" type="text" name="content" placeholder="Write a comment..." ></textarea>
                                    </div>
                                    <div onclick="addComment('${
                                      postData.id
                                    }')" class="postBtn">
                                        <button class="btn btnNotify" type="submit">Add</button>
                                    </div>
                                </form>
                            </div>
                    </div>
                </div>
                `;
        element.innerHTML += x;
      }
      return;
    }
    // if scroll end
    scrollEnd = true;
    element.innerHTML += `<div>There are no posts to show</div>`;
  } catch (error) {
    console.error("Fetching posts failed:", error);
    createAlert(alert("alert-danger", "feald to display post"));
  }
}

async function scrollHandler(api) {
  if (
    window.scrollY + window.innerHeight >=
      document.documentElement.scrollHeight - 20 &&
    Fetching &&
    !scrollEnd
  ) {
    Fetching = false;
    setTimeout(async () => {
      homePage++;
      await display_post(homePage, api);
      Fetching = true;
    }, 500);
  }
}

document.getElementById("centerHome")
  .addEventListener("scroll", () => scrollHandler("filterPost"));
